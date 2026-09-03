package agentbox

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidateIAMTokenName verifies that name can be embedded in the workload-token
// placeholder grammar understood by the AgentBox egress proxy.
func ValidateIAMTokenName(name string) error {
	if name == "" || strings.ContainsAny(name, "{}") || strings.IndexFunc(name, unicode.IsControl) >= 0 {
		return &InvalidArgumentError{Message: fmt.Sprintf("IAM token name %q cannot be empty or contain braces or control characters", name)}
	}
	return nil
}

// IAMTokenPlaceholder returns the value the egress proxy replaces with a
// freshly minted workload token.
func IAMTokenPlaceholder(name string) (string, error) {
	if err := ValidateIAMTokenName(name); err != nil {
		return "", err
	}
	return "${agentbox.identity.tokens." + name + "}", nil
}

// IAMTokenPlaceholders returns placeholders for the supplied registered names.
func IAMTokenPlaceholders(names ...string) (map[string]string, error) {
	result := make(map[string]string, len(names))
	for _, name := range names {
		value, err := IAMTokenPlaceholder(name)
		if err != nil {
			return nil, err
		}
		result[name] = value
	}
	return result, nil
}
