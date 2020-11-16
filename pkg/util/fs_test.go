package util

import (
	"os"
	"path"

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
	It("can create dir if it does not exist", func() {
		tempDir, err := testutil_fs.NewTempDir()
		Expect(err).NotTo(HaveOccurred())
		defer tempDir.Close()

		err = CreateDir(tempDir.Path, 0700)
		Expect(err).NotTo(HaveOccurred())

		tempDir.Close()
		err = CreateDir(tempDir.Path, 0700)
		Expect(err).NotTo(HaveOccurred())

		exists, err := FileExists(tempDir.Path)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())
	})
	It("can create file if it does not exist", func() {
		tempFile, err := testutil_fs.NewTempFile([]byte{})
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		err = CreateFileIfNotExists(tempFile.Path, 0600)
		Expect(err).NotTo(HaveOccurred())
		tempFile.Close()

		tempDir, err := testutil_fs.NewTempDir()
		Expect(err).NotTo(HaveOccurred())
		defer tempDir.Close()

		tempFilePath := path.Join(tempDir.Path, "test")
		err = CreateFileIfNotExists(tempFilePath, 0600)
		Expect(err).NotTo(HaveOccurred())

		exists, err := FileExists(tempFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())

		err = os.Remove(tempFilePath)
		Expect(err).NotTo(HaveOccurred())
	})
	It("can determine homedir", func() {
		Expect(HomeDir()).NotTo(BeEmpty())
	})
})
