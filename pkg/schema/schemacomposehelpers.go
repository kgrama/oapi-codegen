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

func composeTypeDefForSchemaWithPath(path []string, schema *Schema) {
	typeName := PathToTypeName(path)

	typeDef := TypeDefinition{
		TypeName: typeName,
		JsonName: strings.Join(path, "."),
		Schema:   *schema,
	}
	schema.RefType = typeName
	schema.AdditionalTypes = append(schema.AdditionalTypes, typeDef)
}

func translateInteger(f string, outSchema *Schema) {
	// We default to int if format doesn't ask for something else.
	if f == "int64" {
		outSchema.GoType = "int64"
	} else if f == "int32" {
		outSchema.GoType = "int32"
	} else if f == "int16" {
		outSchema.GoType = "int16"
	} else if f == "int8" {
		outSchema.GoType = "int8"
	} else if f == "uint64" {
		outSchema.GoType = "uint64"
	} else if f == "uint32" {
		outSchema.GoType = "uint32"
	} else if f == "uint16" {
		outSchema.GoType = "uint16"
	} else if f == "uint8" {
		outSchema.GoType = "uint8"
	} else if f == "uint" {
		outSchema.GoType = "uint"
	} else {
		outSchema.GoType = "int"
	}
	outSchema.DefineViaAlias = true
}

func translateString(f string, outSchema *Schema) {
	// Special case string formats here.
	switch f {
	case "byte":
		outSchema.GoType = "[]byte"
	case "email":
		outSchema.GoType = "openapi_types.Email"
	case "date":
		outSchema.GoType = "openapi_types.Date"
	case "date-time":
		outSchema.GoType = "time.Time"
	case "json":
		outSchema.GoType = "json.RawMessage"
		outSchema.SkipOptionalPointer = true
	case "uuid":
		outSchema.GoType = "openapi_types.UUID"
	case "binary":
		outSchema.GoType = "openapi_types.File"
	default:
		// All unrecognized formats are simply a regular string.
		outSchema.GoType = "string"
	}
	outSchema.DefineViaAlias = true
}
