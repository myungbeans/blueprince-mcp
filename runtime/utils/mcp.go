package utils

import "fmt"

// ExtractStringParam extracts and validates a string parameter from MCP arguments
func ExtractStringParam(params map[string]any, paramName string) (string, error) {
	if params == nil {
		return "", fmt.Errorf("missing arguments")
	}

	paramRaw, ok := params[paramName]
	if !ok {
		return "", fmt.Errorf("missing required parameter: %s", paramName)
	}

	paramValue, ok := paramRaw.(string)
	if !ok {
		return "", fmt.Errorf("parameter '%s' must be a string", paramName)
	}

	return paramValue, nil
}
