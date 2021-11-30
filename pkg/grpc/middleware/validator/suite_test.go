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

package validator

import (
	"github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	noValidationRules = "No Validation rules on this level"

	validString = "123 Whatever, no re$triction character wise !@#$%^&*()"
	validRestrictedString = "ValidRestricted-String_V1"
	validAddress = "https://k8s-api.lab.example.com:6443"

	validUUID = uuid.New().String()
	validAggregateType = validRestrictedString
	validCSR = []byte("-----BEGIN CERTIFICATE REQUEST-----valid CSR-----END CERTIFICATE REQUEST-----")

	validName = validRestrictedString
	validDisplayName = validString
	validApiServerAddress = validAddress


	invalidStringLength = strings.Repeat("x", 151)
	invalidRestrictedString = "0Start_withNumber-V1"
	invalidRestrictedStringLength = strings.Repeat("x", 61)
	invalidAddress = "not an address"

	invalidUUID = "invalid uuid"
	invalidAggregateTypeStartWithNumber = invalidRestrictedString
	invalidAggregateTypeTooLong = invalidRestrictedStringLength
	invalidCSR = []byte("invalid CSR")

	invalidName = invalidRestrictedString
	invalidDisplayNameTooLong = invalidStringLength
	invalidApiServerAddress = invalidAddress
)

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gRPC Validator Middleware Test Suite")
}

func NewValidCertificateRequest() *commanddata.RequestCertificate {
	return &commanddata.RequestCertificate{
		ReferencedAggregateId:   validUUID,
		ReferencedAggregateType: validAggregateType,
		SigningRequest:          validCSR,
	}
}

func NewValidRequestedCertificate() *eventdata.CertificateRequested {
	return &eventdata.CertificateRequested{
		ReferencedAggregateId:   validUUID,
		ReferencedAggregateType: validAggregateType,
		SigningRequest:          validCSR,
	}
}

func NewValidCreateCluster() *commanddata.CreateCluster {
	return &commanddata.CreateCluster{
		Name: validName,
		DisplayName: validDisplayName,
		ApiServerAddress: validApiServerAddress,
		CaCertBundle: []byte(noValidationRules),
	}
}

func NewValidClusterCreated() *eventdata.ClusterCreated {
	return &eventdata.ClusterCreated{
		Name: validDisplayName,
		Label: validName,
		ApiServerAddress: validApiServerAddress,
		CaCertificateBundle: []byte(noValidationRules),
	}
}

func NewValidClusterCreatedV2() *eventdata.ClusterCreatedV2 {
	return &eventdata.ClusterCreatedV2{
		Name: validName,
		DisplayName: validDisplayName,
		ApiServerAddress: validApiServerAddress,
		CaCertificateBundle: []byte(noValidationRules),
	}
}

func NewValidUpdateCluster() *commanddata.UpdateCluster {
	return &commanddata.UpdateCluster{
		DisplayName: &wrapperspb.StringValue{Value: validDisplayName},
		ApiServerAddress: &wrapperspb.StringValue{Value: validApiServerAddress},
		CaCertBundle: []byte(noValidationRules),
	}
}

func NewValidClusterUpdated() *eventdata.ClusterUpdated {
	return &eventdata.ClusterUpdated{
		DisplayName: validDisplayName,
		ApiServerAddress: validApiServerAddress,
		CaCertificateBundle: []byte(noValidationRules),
	}
}

func NewValidCreateTenantClusterBindingCommandData() *commanddata.CreateTenantClusterBindingCommandData {
	return &commanddata.CreateTenantClusterBindingCommandData{
		TenantId: validUUID,
		ClusterId: validUUID,
	}
}

func NewValidTenantClusterBindingCreated() *eventdata.TenantClusterBindingCreated {
	return &eventdata.TenantClusterBindingCreated{
		TenantId: validUUID,
		ClusterId: validUUID,
	}
}
