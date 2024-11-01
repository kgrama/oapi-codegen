package internal

const (
	// ExtPropGoType overrides the generated type definition.
	ExtPropGoType = "x-go-type"
	// ExtPropGoTypeSkipOptionalPointer specifies that optional fields should
	// be the type itself instead of a pointer to the type.
	ExtPropGoTypeSkipOptionalPointer = "x-go-type-skip-optional-pointer"
	// ExtPropGoImport specifies the module to import which provides above type
	ExtPropGoImport = "x-go-type-import"
	// ExtGoName is used to override a field name
	ExtGoName = "x-go-name"
	// ExtGoTypeName is used to override a generated typename for something.
	ExtGoTypeName        = "x-go-type-name"
	ExtPropGoJsonIgnore  = "x-go-json-ignore"
	ExtPropOmitEmpty     = "x-omitempty"
	ExtPropExtraTags     = "x-oapi-codegen-extra-tags"
	ExtEnumVarNames      = "x-enum-varnames"
	ExtEnumNames         = "x-enumNames"
	ExtDeprecationReason = "x-deprecated-reason"
	ExtOrder             = "x-order"
	// ExtOapiCodegenOnlyHonourGoName is to be used to explicitly enforce the generation of a field as the `x-go-name` extension has describe it.
	// This is intended to be used alongside the `allow-unexported-struct-field-names` Compatibility option
	ExtOapiCodegenOnlyHonourGoName = "x-oapi-codegen-only-honour-go-name"
)
