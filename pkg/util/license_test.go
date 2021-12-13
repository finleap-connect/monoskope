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

package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var skipPaths = []string{
	"../../pkg/api/",
	"../../.go/",
	"../../test/",
	"../../include/",
}

var _ = Describe("License", func() {
	It("should exist in every .go-file", func() {
		err := filepath.Walk("../..", func(path string, fi os.FileInfo, err error) error {
			for _, skipPath := range skipPaths {
				if strings.HasPrefix(path, skipPath) {
					return nil
				}
			}
			if err != nil {
				return err
			}
			if filepath.Ext(path) != ".go" {
				return nil
			}
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return nil
			}
			if !bytes.HasPrefix(content, []byte("// Copyright")) {
				return fmt.Errorf("%s: license header missing", path)
			}
			return nil
		})
		Expect(err).ToNot(HaveOccurred())
	})
})
