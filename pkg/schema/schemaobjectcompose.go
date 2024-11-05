package schema

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/globalstate"
)

func composeGoSchemaForObject(schema *openapi3.Schema, t *openapi3.Types, outSchema *Schema, path []string) (Schema, error) {
	if len(schema.Properties) == 0 && !SchemaHasAdditionalProperties(schema) && schema.AnyOf == nil && schema.OneOf == nil {
		setSchemaoutputTypeMaporNone(outSchema, t)
	} else {
		setOtherPropsForEmptyOrObj(outSchema, schema)

		// If additional properties are defined, we will override the default
		// above with the specific definition.
		if schema.AdditionalProperties.Schema != nil {
			additionalSchema, err := GenerateGoSchema(schema.AdditionalProperties.Schema, path)
			if err != nil {
				return Schema{}, fmt.Errorf("error generating type for additional properties: %w", err)
			}
			if additionalSchema.HasAdditionalProperties || len(additionalSchema.UnionElements) != 0 {
				// If we have fields present which have additional properties or union values,
				// but are not a pre-defined type, we need to define a type
				// for them, which will be based on the field names we followed
				// to get to the type.
				composeTypeDefForSchemaWithPath(append(path, "AdditionalProperties"), &additionalSchema)
			}
			outSchema.AdditionalPropertiesType = &additionalSchema
			outSchema.AdditionalTypes = append(outSchema.AdditionalTypes, additionalSchema.AdditionalTypes...)
		}

		// If the schema has no properties, and only additional properties, we will
		// early-out here and generate a map[string]<schema> instead of an object
		// that contains this map. We skip over anyOf/oneOf here because they can
		// introduce properties. allOf was handled above.
		if !globalstate.GlobalState.Options.Compatibility.DisableFlattenAdditionalProperties &&
			len(schema.Properties) == 0 && schema.AnyOf == nil && schema.OneOf == nil {
			// We have a dictionary here. Returns the goType to be just a map from
			// string to the property type. HasAdditionalProperties=false means
			// that we won't generate custom json.Marshaler and json.Unmarshaler functions,
			// since we don't need them for a simple map.
			outSchema.HasAdditionalProperties = false
			outSchema.GoType = fmt.Sprintf("map[string]%s", additionalPropertiesType(*outSchema))
			return *outSchema, nil
		}

		// We've got an object with some properties.
		for _, pName := range SortedSchemaKeys(schema.Properties) {
			p := schema.Properties[pName]
			propertyPath := append(path, pName)
			pSchema, err := GenerateGoSchema(p, propertyPath)
			if err != nil {
				return Schema{}, fmt.Errorf("error generating Go schema for property '%s': %w", pName, err)
			}

			required := StringInArray(pName, schema.Required)

			if (pSchema.HasAdditionalProperties || len(pSchema.UnionElements) != 0) && pSchema.RefType == "" {
				// If we have fields present which have additional properties or union values,
				// but are not a pre-defined type, we need to define a type
				// for them, which will be based on the field names we followed
				// to get to the type.
				composeTypeDefForSchemaWithPath(propertyPath, &pSchema)
			}
			description := ""
			if p.Value != nil {
				description = p.Value.Description
			}
			prop := Property{
				JsonFieldName: pName,
				Schema:        pSchema,
				Required:      required,
				Description:   description,
				Nullable:      p.Value.Nullable,
				ReadOnly:      p.Value.ReadOnly,
				WriteOnly:     p.Value.WriteOnly,
				Extensions:    p.Value.Extensions,
				Deprecated:    p.Value.Deprecated,
			}
			outSchema.Properties = append(outSchema.Properties, prop)
			if len(pSchema.AdditionalTypes) > 0 {
				outSchema.AdditionalTypes = append(outSchema.AdditionalTypes, pSchema.AdditionalTypes...)
			}
		}

		if schema.AnyOf != nil {
			if err := generateUnion(outSchema, schema.AnyOf, schema.Discriminator, path); err != nil {
				return Schema{}, fmt.Errorf("error generating type for anyOf: %w", err)
			}
		}
		if schema.OneOf != nil {
			if err := generateUnion(outSchema, schema.OneOf, schema.Discriminator, path); err != nil {
				return Schema{}, fmt.Errorf("error generating type for oneOf: %w", err)
			}
		}

		outSchema.GoType = GenStructFromSchema(*outSchema)
	}

	// Check for x-go-type-name. It behaves much like x-go-type, however, it will
	// create a type definition for the named type, and use the named type in place
	// of this schema.
	if extension, ok := schema.Extensions[extGoTypeName]; ok {
		typeName, err := extTypeName(extension)
		if err != nil {
			return *outSchema, fmt.Errorf("invalid value for %q: %w", extGoTypeName, err)
		}

		newTypeDef := TypeDefinition{
			TypeName: typeName,
			Schema:   *outSchema,
		}
		outSchema = &Schema{
			Description:     newTypeDef.Schema.Description,
			GoType:          typeName,
			DefineViaAlias:  true,
			AdditionalTypes: append(outSchema.AdditionalTypes, newTypeDef),
		}
	}

	return *outSchema, nil
}
