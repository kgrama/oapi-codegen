package schema

import (
	"strings"
)

// UnionElement describe union element, based on prefix externalRef\d+ and real ref name from external schema.
type UnionElement string

// String returns externalRef\d+ and real ref name from external schema, like externalRef0.SomeType.
func (u UnionElement) String() string {
	return string(u)
}

// Method generate union method name for template functions `As/From/Merge`.
func (u UnionElement) Method() string {
	var method string
	for _, part := range strings.Split(string(u), `.`) {
		method += UppercaseFirstCharacter(part)
	}
	return method
}
