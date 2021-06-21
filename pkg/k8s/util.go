package k8s

import (
	"errors"
	"regexp"
	"strings"
)

// GetNamespaceName returns a sanitized version of the input string
// which is allowed to be used as a kubernetes namespace.
func GetNamespaceName(any string) (string, error) {
	nsName := strings.ToLower(any)
	replacer := strings.NewReplacer(" ", "",
		"ü", "ue",
		"ö", "oe",
		"ä", "ae",
		"ß", "ss",
		"_", "-",
		".", "-",
		"/", "-")
	nsName = replacer.Replace(nsName)

	// regex for checking k8s namespace name
	regex, err := regexp.Compile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	if err != nil {
		return "", err
	}

	if !regex.MatchString(nsName) {
		return "", errors.New("namespace name does not adhere to the naming rules")
	}
	return nsName, nil
}
