package openapiv3

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func FilterOperationsByTag(spec *openapi3.T, opts Configuration) {
	filterOperationsByTag(spec, opts)
}

func FilterOperationsByOperationID(spec *openapi3.T, opts Configuration) {
	filterOperationsByOperationID(spec, opts)
}

func PruneUnusedComponents(spec *openapi3.T) {
	pruneUnusedComponents(spec)
}
