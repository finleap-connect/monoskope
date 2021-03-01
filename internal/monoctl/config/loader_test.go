package config

import (
	"os"
	"time"

	testutil_fs "github.com/kubism/testutil/pkg/fs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	fakeConfigData = `server: https://1.1.1.1`
)

var _ = Describe("client config loader", func() {
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
		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		loader := NewLoaderFromExplicitFile(tempFile.Path)
		err = loader.LoadAndStoreConfig()

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
	})
	It("can init config for explicit file path", func() {
		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		conf := NewConfig()
		conf.Server = "localhost"

		loader := NewLoaderFromExplicitFile(tempFile.Path)
		os.Remove(tempFile.Path)
		err = loader.InitConifg(conf, false)

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
	})
	It("can init config for env var path", func() {
		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		loader := NewLoader()
		conf := NewConfig()
		conf.Server = "localhost"

		os.Setenv(RecommendedConfigPathEnvVar, tempFile.Path)
		os.Remove(tempFile.Path)
		err = loader.InitConifg(conf, false)

		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
	})
	It("can save config", func() {
		tempFile, err := testutil_fs.NewTempFile([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		defer tempFile.Close()

		loader := NewLoaderFromExplicitFile(tempFile.Path)
		err = loader.LoadAndStoreConfig()
		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())

		conf := loader.GetConfig()
		conf.Server = "monoskope.io"
		err = loader.SaveConfig()
		Expect(err).NotTo(HaveOccurred())

		loader = NewLoaderFromExplicitFile(tempFile.Path)
		err = loader.LoadAndStoreConfig()
		Expect(err).NotTo(HaveOccurred())
		Expect(loader.config).ToNot(BeNil())
		Expect(loader.config.Server).To(Equal("monoskope.io"))
	})
	It("can validate", func() {
		loader := NewLoader()
		conf, err := loader.LoadFromBytes([]byte(fakeConfigData))
		Expect(err).NotTo(HaveOccurred())
		Expect(conf).ToNot(BeNil())

		Expect(conf.HasAuthInformation()).To(BeFalse())

		conf.AuthInformation = &AuthInformation{}
		Expect(conf.HasAuthInformation()).To(BeTrue())
		Expect(conf.AuthInformation.HasToken()).To(BeFalse())
		Expect(conf.AuthInformation.HasRefreshToken()).To(BeFalse())
		Expect(conf.AuthInformation.IsValid()).To(BeFalse())

		conf.AuthInformation.Token = "test"
		Expect(conf.AuthInformation.HasToken()).To(BeTrue())
		Expect(conf.AuthInformation.HasRefreshToken()).To(BeFalse())
		Expect(conf.AuthInformation.IsValid()).To(BeTrue())

		conf.AuthInformation.RefreshToken = "testrefresh"
		Expect(conf.AuthInformation.HasToken()).To(BeTrue())
		Expect(conf.AuthInformation.HasRefreshToken()).To(BeTrue())
		Expect(conf.AuthInformation.IsValid()).To(BeTrue())

		expiry := time.Now().Add(1 * time.Hour)
		conf.AuthInformation.Expiry = &expiry
		Expect(conf.AuthInformation.HasToken()).To(BeTrue())
		Expect(conf.AuthInformation.HasRefreshToken()).To(BeTrue())
		Expect(conf.AuthInformation.IsTokenExpired()).To(BeFalse())
		Expect(conf.AuthInformation.IsValid()).To(BeTrue())

		expiry = time.Now().Add(-1 * time.Hour)
		conf.AuthInformation.Expiry = &expiry
		Expect(conf.AuthInformation.HasToken()).To(BeTrue())
		Expect(conf.AuthInformation.HasRefreshToken()).To(BeTrue())
		Expect(conf.AuthInformation.IsTokenExpired()).To(BeTrue())
		Expect(conf.AuthInformation.IsValid()).To(BeFalse())
	})
})
