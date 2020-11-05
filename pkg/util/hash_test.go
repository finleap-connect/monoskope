package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("util.hash", func() {
	It("can hash a string", func() {
		testString := "this is a test"
		expectedResult := "2e99758548972a8e8822ad47fa1017ff72f06f3ff6a016851f45c398732bc50c"
		hashedString := HashString(testString)
		Expect(hashedString).To(Equal(expectedResult))
	})
})
