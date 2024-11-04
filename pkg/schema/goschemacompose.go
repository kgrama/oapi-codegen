package schema

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func refToAnother(schema *openapi3.Schema, ref string) (Schema, error) {

	// Convert the reference path to Go type
	refType, err := RefPathToGoType(ref)
	if err != nil {
		return Schema{}, fmt.Errorf("error turning reference (%s) into a Go type: %s",
			ref, err)
	}
	return Schema{
		GoType:         refType,
		Description:    schema.Description,
		DefineViaAlias: true,
		OAPISchema:     schema,
	}, nil
}

func applySchemaMerge(schema *openapi3.Schema, path []string) (Schema, error) {
	mergedSchema, err := MergeSchemas(schema.AllOf, path)
	if err != nil {
		return Schema{}, fmt.Errorf("error merging schemas: %w", err)
	}
	mergedSchema.OAPISchema = schema
	return mergedSchema, nil
}

func overrideSchema(extension interface{}, schema *Schema) error {
	typeName, err := extTypeName(extension)
	if err != nil {
		return fmt.Errorf("invalid value for %q: %w", extPropGoType, err)
	}
	schema.GoType = typeName
	schema.DefineViaAlias = true
	return nil
}

func optionalPointer(extension interface{}, schema *Schema) error {
	skipOptionalPointer, err := extParsePropGoTypeSkipOptionalPointer(extension)
	if err != nil {
		return fmt.Errorf("invalid value for %q: %w", extPropGoTypeSkipOptionalPointer, err)
	}
	schema.SkipOptionalPointer = skipOptionalPointer
	return nil
}

func setSchemaoutputTypeMaporNone(outSchema *Schema, t *openapi3.Types) {
	// If the object has no properties or additional properties, we
	// have some special cases for its type.
	if t.Is("object") {
		// We have an object with no properties. This is a generic object
		// expressed as a map.
		outSchema.GoType = "map[string]interface{}"
	} else { // t == ""
		// If we don't even have the object designator, we're a completely
		// generic type.
		outSchema.GoType = "interface{}"
	}
	outSchema.DefineViaAlias = true
}

func setOtherPropsForEmptyOrObj(outSchema *Schema, schema *openapi3.Schema) {
	// When we define an object, we want it to be a type definition,
	// not a type alias, eg, "type Foo struct {...}"
	outSchema.DefineViaAlias = false

	// If the schema has additional properties, we need to special case
	// a lot of behaviors.
	outSchema.HasAdditionalProperties = SchemaHasAdditionalProperties(schema)

	// Until we have a concrete additional properties type, we default to
	// any schema.
	outSchema.AdditionalPropertiesType = &Schema{
		GoType: "interface{}",
	}
}

func composeTypedefForSchemaandPath(path []string, schema *Schema) {
	typeName := PathToTypeName(path)

	typeDef := TypeDefinition{
		TypeName: typeName,
		JsonName: strings.Join(path, "."),
		Schema:   *schema,
	}
	schema.RefType = typeName
	schema.AdditionalTypes = append(schema.AdditionalTypes, typeDef)
}
