package schema

import (
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/globalstate"
)

type Property struct {
	Description   string
	JsonFieldName string
	Schema        Schema
	Required      bool
	Nullable      bool
	ReadOnly      bool
	WriteOnly     bool
	NeedsFormTag  bool
	Extensions    map[string]interface{}
	Deprecated    bool
}

func (p Property) GoFieldName() string {
	goFieldName := p.JsonFieldName
	if extension, ok := p.Extensions[extGoName]; ok {
		if extGoFieldName, err := extParseGoFieldName(extension); err == nil {
			goFieldName = extGoFieldName
		}
	}

	if globalstate.GlobalState.Options.Compatibility.AllowUnexportedStructFieldNames {
		if extension, ok := p.Extensions[extOapiCodegenOnlyHonourGoName]; ok {
			if extOapiCodegenOnlyHonourGoName, err := extParseOapiCodegenOnlyHonourGoName(extension); err == nil {
				if extOapiCodegenOnlyHonourGoName {
					return goFieldName
				}
			}
		}
	}

	return SchemaNameToTypeName(goFieldName)
}

func (p Property) GoTypeDef() string {
	typeDef := p.Schema.TypeDecl()
	if globalstate.GlobalState.Options.OutputOptions.NullableType && p.Nullable {
		return "nullable.Nullable[" + typeDef + "]"
	}
	if !p.Schema.SkipOptionalPointer &&
		(!p.Required || p.Nullable ||
			(p.ReadOnly && (!p.Required || !globalstate.GlobalState.Options.Compatibility.DisableRequiredReadOnlyAsPointer)) ||
			p.WriteOnly) {

		typeDef = "*" + typeDef
	}
	return typeDef
}

func PropertiesEqual(a, b Property) bool {
	return a.JsonFieldName == b.JsonFieldName && a.Schema.TypeDecl() == b.Schema.TypeDecl() && a.Required == b.Required
}
