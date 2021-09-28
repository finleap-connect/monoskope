// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
