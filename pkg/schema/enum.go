package schema

// EnumDefinition holds type information for enum
type EnumDefinition struct {
	// Schema is the scheme of a type which has a list of enum values, eg, the
	// "container" of the enum.
	Schema Schema
	// TypeName is the name of the enum's type, usually aliased from something.
	TypeName string
	// ValueWrapper wraps the value. It's used to conditionally apply quotes
	// around strings.
	ValueWrapper string
	// PrefixTypeName determines if the enum value is prefixed with its TypeName.
	// This is set to true when this enum conflicts with another in terms of
	// TypeNames or when explicitly requested via the
	// `compatibility.always-prefix-enum-values` option.
	PrefixTypeName bool
}

// GetValues generates enum names in a way to minimize global conflicts
func (e *EnumDefinition) GetValues() map[string]string {
	// in case there are no conflicts, it's safe to use the values as-is
	if !e.PrefixTypeName {
		return e.Schema.EnumValues
	}
	// If we do have conflicts, we will prefix the enum's typename to the values.
	newValues := make(map[string]string, len(e.Schema.EnumValues))
	for k, v := range e.Schema.EnumValues {
		newName := e.TypeName + UppercaseFirstCharacter(k)
		newValues[newName] = v
	}
	return newValues
}
