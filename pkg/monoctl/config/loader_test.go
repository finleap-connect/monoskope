package config

import (
	"os"

	testutil_fs "github.com/kubism/testutil/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	fakeConfigData = `server: https://1.1.1.1`
)

var _ = Describe("loader", func() {
	It("can load config from bytes", func() {
		loader := NewLoader()
		conf, err := loader.LoadFromBytes([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		Expect(conf).ToNot(BeNil())
	})
	It("errors for empty config", func() {
		loader := NewLoader()
		conf, err := loader.LoadFromBytes([]byte(""))
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ErrEmptyServer))
		Expect(conf).To(BeNil())
	})
	It("loads config from env var path", func() {
		loader := NewLoader()

		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		os.Setenv(RecommendedConfigPathEnvVar, tempFile.Path)
		err = loader.LoadAndStoreConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
	})
	It("loads config from explicit file path", func() {
		loader := NewLoader()

		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		loader.ExplicitFile = tempFile.Path
		err = loader.LoadAndStoreConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
	})
})
