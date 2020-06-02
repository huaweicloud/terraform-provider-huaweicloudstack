package huaweicloudstack

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func ValidateStringList(v interface{}, k string, l []string) (ws []string, errors []error) {
	value := v.(string)
	for i := range l {
		if value == l[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, l))
	return
}

// Not currently used.
func ValidateInsensitiveStringList(v interface{}, k string, l []string) (ws []string, errors []error) {
	value := v.(string)
	for i := range l {
		if strings.EqualFold(value, l[i]) {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, l))
	return
}

func ValidateIntRange(v interface{}, k string, l int, h int) (ws []string, errors []error) {
	i, ok := v.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("%q must be an integer", k))
		return
	}
	if i < l || i > h {
		errors = append(errors, fmt.Errorf("%q must be between %d and %d", k, l, h))
		return
	}
	return
}

func validateTrueOnly(v interface{}, k string) (ws []string, errors []error) {
	if b, ok := v.(bool); ok && b {
		return
	}
	if v, ok := v.(string); ok && v == "true" {
		return
	}
	errors = append(errors, fmt.Errorf("%q must be true", k))
	return
}

func validateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

func validateKmsKeyStatus(v interface{}, k string) (ws []string, errors []error) {
	status := v.(string)
	if status != EnabledState && status != DisabledState && status != PendingDeletionState {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid status, expected %s or %s or %s, got %s.",
			k, EnabledState, DisabledState, PendingDeletionState, status))
	}
	return
}

func validateName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}

	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

func validateIP(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ipnet := net.ParseIP(value)

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network IP address, got %q", k, value))
	}

	return
}

func validateCIDR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, got %q", k, value))
	}

	return
}

func validateStackTemplate(v interface{}, k string) (ws []string, errors []error) {
	if looksLikeJsonString(v) {
		if _, err := normalizeJsonString(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
		}
	} else {
		if _, err := checkYamlString(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %s", k, err))
		}
	}
	return
}
