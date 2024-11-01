package codegen

import (
	"fmt"
)

func extString(extPropValue interface{}) (string, error) {
	str, ok := extPropValue.(string)
	if !ok {
		return "", fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	return str, nil
}

func extTypeName(extPropValue interface{}) (string, error) {
	return extString(extPropValue)
}

func extParsePropGoTypeSkipOptionalPointer(extPropValue interface{}) (bool, error) {
	goTypeSkipOptionalPointer, ok := extPropValue.(bool)
	if !ok {
		return false, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	return goTypeSkipOptionalPointer, nil
}

func extParseGoFieldName(extPropValue interface{}) (string, error) {
	return extString(extPropValue)
}

func extParseOmitEmpty(extPropValue interface{}) (bool, error) {
	omitEmpty, ok := extPropValue.(bool)
	if !ok {
		return false, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	return omitEmpty, nil
}

func extExtraTags(extPropValue interface{}) (map[string]string, error) {
	tagsI, ok := extPropValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	tags := make(map[string]string, len(tagsI))
	for k, v := range tagsI {
		vs, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert type: %T", v)
		}
		tags[k] = vs
	}
	return tags, nil
}

func extParseGoJsonIgnore(extPropValue interface{}) (bool, error) {
	goJsonIgnore, ok := extPropValue.(bool)
	if !ok {
		return false, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	return goJsonIgnore, nil
}

func extParseEnumVarNames(extPropValue interface{}) ([]string, error) {
	namesI, ok := extPropValue.([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	names := make([]string, len(namesI))
	for i, v := range namesI {
		vs, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert type: %T", v)
		}
		names[i] = vs
	}
	return names, nil
}

func extParseDeprecationReason(extPropValue interface{}) (string, error) {
	return extString(extPropValue)
}

func extParseOapiCodegenOnlyHonourGoName(extPropValue interface{}) (bool, error) {
	onlyHonourGoName, ok := extPropValue.(bool)
	if !ok {
		return false, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	return onlyHonourGoName, nil
}
