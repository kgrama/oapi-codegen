package schema

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Schema describes an OpenAPI schema, with lots of helper fields to use in the
// templating engine.
type Schema struct {
	GoType  string // The Go type needed to represent the schema
	RefType string // If the type has a type name, this is set

	ArrayType *Schema // The schema of array element

	EnumValues map[string]string // Enum values

	Properties               []Property       // For an object, the fields with names
	HasAdditionalProperties  bool             // Whether we support additional properties
	AdditionalPropertiesType *Schema          // And if we do, their type
	AdditionalTypes          []TypeDefinition // We may need to generate auxiliary helper types, stored here

	SkipOptionalPointer bool // Some types don't need a * in front when they're optional

	Description string // The description of the element

	UnionElements []UnionElement // Possible elements of oneOf/anyOf union
	Discriminator *Discriminator // Describes which value is stored in a union

	// If this is set, the schema will declare a type via alias, eg,
	// `type Foo = bool`. If this is not set, we will define this type via
	// type definition `type Foo bool`
	//
	// Can be overriden by the OutputOptions#DisableTypeAliasesForType field
	DefineViaAlias bool

	// The original OpenAPIv3 Schema.
	OAPISchema *openapi3.Schema
}

func (s Schema) IsRef() bool {
	return s.RefType != ""
}

func (s Schema) IsExternalRef() bool {
	if !s.IsRef() {
		return false
	}
	return strings.Contains(s.RefType, ".")
}

func (s Schema) TypeDecl() string {
	if s.IsRef() {
		return s.RefType
	}
	return s.GoType
}

// AddProperty adds a new property to the current Schema, and returns an error
// if it collides. Two identical fields will not collide, but two properties by
// the same name, but different definition, will collide. It's safe to merge the
// fields of two schemas with overlapping properties if those properties are
// identical.
func (s *Schema) AddProperty(p Property) error {
	// Scan all existing properties for a conflict
	for _, e := range s.Properties {
		if e.JsonFieldName == p.JsonFieldName && !PropertiesEqual(e, p) {
			return fmt.Errorf("property '%s' already exists with a different type", e.JsonFieldName)
		}
	}
	s.Properties = append(s.Properties, p)
	return nil
}

func (s Schema) GetAdditionalTypeDefs() []TypeDefinition {
	return s.AdditionalTypes
}
