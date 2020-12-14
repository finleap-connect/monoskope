package util

import (
	testutil_fs "github.com/kubism/testutil/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("util.fs", func() {
	It("can check if file exists", func() {
		tempDir, err := testutil_fs.NewTempDir()
		Expect(err).NotTo(HaveOccurred())
		defer tempDir.Close()

		exists, err := FileExists(tempDir.Path)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())

		tempDir.Close()
		exists, err = FileExists(tempDir.Path)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})
	It("can determine homedir", func() {
		Expect(HomeDir()).NotTo(BeEmpty())
	})
})
