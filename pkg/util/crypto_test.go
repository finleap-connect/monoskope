// Copyright 2022 Monoskope Authors
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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("util.EncryptAES", func() {
	plaintextValue := "this is my secret"
	encryptionKey := "thisis32bitlongpassphraseimusing"
	var encryptedBytes []byte

	It("can encrypt bytes", func() {
		var err error
		encryptedBytes, err = EncryptAES([]byte(encryptionKey), []byte(plaintextValue))
		Expect(err).ToNot(HaveOccurred())
		Expect(encryptedBytes).ToNot(BeNil())
	})
	It("can decrypt bytes", func() {
		decryptedBytes, err := DecryptAES([]byte(encryptionKey), encryptedBytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(decryptedBytes).ToNot(BeNil())
		Expect(string(decryptedBytes)).To(Equal(plaintextValue))
	})
})
