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

package k8sauthz

import (
	_ "embed"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:embed test_config.yaml
var test_config []byte

var _ = Describe("internal/k8sauthz", func() {
	Context("GitRepository", func() {
		It("NewGitRepository() creates a new instance with defaults", func() {
			os.Setenv("test1.basic.username", "test1")
			os.Setenv("test1.basic.password", "testpw")

			os.Setenv("test2.ssh.privateKey", "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcGdJQkFBS0NBUUVBeGdxM0E2Y2dBWEEwM2g4SXY4cWxLdnZpeHhLT0FNcndwYTBRazBpTnJOVzdzOTJ4CnhmQ1BzQjJpSlBnOFRLNGpoSUJyRlhSNWNrODBoSitwb0ltWVpQM2cyYXlUUEpBRU9JWVVQN3I0UjNOYXJGTXMKbU9xWGliQURoUGFRaFZZYkhSUEw1eHFlKzBselhpQk5rQVJWTUo4VkFqMUZkZkZXOTAwbWllQkl0VUxxdUdnUwo0UFplWE81eTlhNW9PZWQ1bFBnNmJOWG5UMWYxWDVxaktISWw1MXc4M3JhSHgzd1l6YUJaeXRKejNNVjhFQWFyClNwVnQvdFR3MnJha2Y1YVFJZUI2ZitFV3BZRUpBK3ZFQVFWaUlaTEl4Ly90T0RkU3ZRYjFLakhzUnp4UEJDWjcKU1VOVzFkT3AxSTdjdEtpaWdlSmRXQ3grZnFWblYralRGM2NsN1FJREFRQUJBb0lCQVFEQURMQUNDTllPendOVgp5LzZ1RHhReEpPZDhhYy92a04zaHJIMEFkMEY4dENBOGd1YmNyemFJY1pSR0NieGdHQmMxYlZ0ckNQS2xHWTR5CmRxc0dsYWlGZDNYUFlYMi9JMklVZG1HNGs5WWxaQVl3U3RCYTBsUjJINHlxS25sbjlHS0N4aW5jU1lLVzZWbkoKazhYanZXL25vTDI4MkRJOTkvYUFCLzh4Z2FEQXRSWmZORS9WaVZYZm1XN3JRdTZIMmFWY3VqakxXbURRbFJ5cwpDcTgzc3I1QXpvUWYxdWwvaTkyTXZ6eUN2cVFreHI1NndLMjREU0FmU3FNMmxQdE5oL3BRQjlGTWtZSGUzYWlSCkJjNWorRndMZ1g0UHZHdUNXc0xMajU4Q3haWmkwTVBTUmhYRmhxVVhNd3gwV2lFK2tZSFRNUElsMjk4b0hBbnYKWkNVa1FOcTVBb0dCQVB4bDdPaTRMQ1pveGh2YkUxMnE0eitYUngvTFJOZFVtdy9tdDFJUUhXUlpjOHIya3lkYwpnMllUS2FHdSszeUpKUHcwRmlsdzFQYnlyU2x3RVpBNzY4RE84SDFqajJQVFlpaExGRm1UNzFPbXRlZ0VXSlEyClVQNm1NcTVBQzdPaFVGYTVYYi9OcUdRSW9qR0lHVVBYREQweVBqbkhGNzhsSlpMRytGR1JWYmgvQW9HQkFNamUKTmszK1lLOThxL09MR0JqNTJ6bFVNVGtIVmN6L0wyYThoWmc0bnNpK0FoTnFOVjVndC9SZ0dsZW5FOUdSNDZpUAprQVpsTlZHVG9laDFQT1hKVUQ1TUd2Wit0cTR6L3g5dnRORWtqODhZUWV0N1ozOXNqY3QrK2pIcEk4N09ENVZoCitGWCtQZklrUFVIbjhNYnp1TlVVUmJrL0lEY2V6amRtdFhmcFQwdVRBb0dCQU5FdzhWQlpCVTY0WENwT1F4akoKUUJ3K2w5YUVOUVI3dlNGS3lmb2NIU0JFKzdIbStFUVJhMTY0MXpLSXd5dU94N2E0dlAyUDVSTGdSQ0VxSDRSNgpCWVlDZmxTT3BoUEk2WHdYN2ovb1d0M3NPS3lhbllnS0ZNdGFtSHJQM21MOWVvYkdrQ3NlbTBoOTEyQlNPdzU5Ck9FbW9VT2EyV3JvMFlaWHJqM2liaW42bEFvR0JBTEZ5M1RaNWlxQjgyc3NFRGYzQ2ZOQmdlRHVSSjFNNS9INGQKL1VkRWZjR0pXZndjejVqWWlLbnlYRk1pM25jOUVvbE9pa1djRTBaRnpicTJGMTVJUWdORzZHcnA5aWhlOXZxbQpqRzVXaGxURmNUYStoZWdqMWYvMzVMOUMzc2RMY3FqZEs2Mk1OTjA4OW9ES2pnSzBQSXpBby9mS3RJTDlTOS9WCmRHckNTYkxSQW9HQkFMRW5zNXRLY2h2dFNVZ0ZjRDB3a0Jvdnd4TWV3TklkSjFNN0F0RTRmQnpYN0cya0tNSDUKVUpZK0ErUFZ2cTA1RnF4S3FGY0JSTFQrZW05RkdnVUNlbVZCYlg2S1R1eUxkOHRBRnpaWm9tUm9TQTFVbW05QQpxamJweHkvNTMyRlpYdjQrSzNaN1BtSnFUVUtoZ09uTS9VREh6YXBCUlhHdHlpY0lvcU9FaUJtaQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=")
			os.Setenv("test2.ssh.password", "a25vd24taG9zdC1rZXlz")

			conf, err := NewConfigFromFile(test_config)
			Expect(err).ToNot(HaveOccurred())
			Expect(conf).ToNot(BeNil())
			Expect(len(conf.Repositories)).To(BeNumerically("==", 2))
			Expect(len(conf.Mappings)).To(BeNumerically("==", 2))

			firstRepo := conf.Repositories[0]
			Expect(firstRepo.AllClusters).To(BeTrue())
			Expect(len(firstRepo.Clusters)).To(BeNumerically("==", 0))

			secondRepo := conf.Repositories[1]
			Expect(secondRepo.AllClusters).To(BeFalse())
			Expect(len(secondRepo.Clusters)).To(BeNumerically("==", 2))
		})
	})
})
