package schema

import (
	"fmt"
	"strings"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/globalstate"
)

// ensureExternalRefsInSchema ensures that when an externalRef (`$ref` that points to a file that isn't the current spec) is encountered, we make sure we update our underlying `RefType` to make sure that we point to that type.
//
// This only happens if we have a non-empty `ref` passed in, and that `ref` isn't pointing to something in our file
//
// NOTE that the pointer here allows us to pass in a reference and edit in-place
func ensureExternalRefsInSchema(schema *Schema, ref string) {
	if ref == "" {
		return
	}

	// if this is already defined as the start of a struct, we shouldn't inject **??**
	if strings.HasPrefix(schema.GoType, "struct {") {
		return
	}

	parts := strings.SplitN(ref, "#", 2)
	if pack, ok := globalstate.GlobalState.ImportMapping[parts[0]]; ok {
		schema.RefType = fmt.Sprintf("%s.%s", pack.Name, schema.GoType)
	}
}
