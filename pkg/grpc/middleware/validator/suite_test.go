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

package validator

import (
	"strings"
	"testing"

	"github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing/commands"
	"github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	validString           = "123 Whatever, no re$triction character wise !@#$%^&*()"
	validRestrictedString = "ValidRestricted-String_V1"
	validLowercaseString  = "onlylowercase"

	validUUID          = uuid.New().String()
	validAggregateType = validRestrictedString

	validName             = validRestrictedString
	validDisplayName      = validString
	validApiServerAddress = "https://k8s-api.lab.example.com:6443"

	validTenantPrefix = validRestrictedString[0:12]

	validEmail = "email@valid.com"
	validRole  = validLowercaseString
	validScope = validLowercaseString

	validCommandType = validRestrictedString

	validEventType = validRestrictedString

	validCertificate = []byte(`-----BEGIN CERTIFICATE-----
	MIICnTCCAkSgAwIBAgIQMo7x823NtJ/Xyy1Wl+8+yzAKBggqhkjOPQQDAjAnMSUw
	IwYDVQQDExxyb290Lm1vbm9za29wZS5jbHVzdGVyLmxvY2FsMB4XDTIxMDYwMjAy
	MTAxNVoXDTIxMDYwNDAyMTAxNVowMDESMBAGA1UEChMJTW9ub3Nrb3BlMRowGAYD
	VQQDExFtOC1hdXRoZW50aWNhdGlvbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
	AQoCggEBANCKZWW0el3OzPw7914TC1Ld2At/xIh/3zoiawQcbS8mrjnVMO2oSomY
	mks6sEaWp4p80PwJkzSplpgoJmEOYqps+YXo+1NLp66bFPkAbMEZDsZ4QmrQQ7X3
	iv5IaDFW4vSGJFSkTQnUmedlhrWguasOD3vL0Pek89L8kQ09+YlDk/fpBZUXFADU
	+ef4GjTkWJzkg32dSOudJDYD4wUPczTFlRO097MBBlaMb4LKYfDfjuUKRCOAL3LD
	7kKAatHKeoADuBptUv/lQLExGNzlhRteaLocTHHab2hs+NCFYABv2Px5Tcnbw8g+
	/r/97gwKkpFeF5p4WhdVgbDYd2MGUlMCAwEAAaN+MHwwHQYDVR0lBBYwFAYIKwYB
	BQUHAwIGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAUI+RyVqj0
	J9qH8l3pbY9KUkoTgHQwLAYDVR0RBCUwI4IJbG9jYWxob3N0hwR/AAABhxAAAAAA
	AAAAAAAAAAAAAAAAMAoGCCqGSM49BAMCA0cAMEQCIEPbvMo2YvqlYQtdkQwlhJci
	mTlsDv6VmO4WfCjrQdwLAiA+N0eeiL/yLPC5ReaPYQ7PeoXbc9+EPR2FBDrkiBbA
	8w==
	-----END CERTIFICATE-----`)

	invalidStringLength     = strings.Repeat("x", 151)
	invalidRestrictedString = "0Start_withNumber-V1"
	invalidLowercaseString  = "onlyLowerCase"
	invalidStringWhitespace = " " + validString + "\n"

	invalidUUID                         = "invalid uuid"
	invalidAggregateTypeStartWithNumber = invalidRestrictedString

	invalidName                   = invalidRestrictedString
	invalidDisplayNameTooLong     = invalidStringLength
	invalidDisplayNameWhiteSpaces = invalidStringWhitespace
	invalidApiServerAddress       = "k8s-api.lab. example.com:6443"

	invalidTenantPrefixTooLong         = validRestrictedString
	invalidTenantPrefixStartWithNumber = invalidRestrictedString[0:12]

	invalidEmail = "email#invalid.com"
	invalidRole  = invalidLowercaseString
	invalidScope = invalidLowercaseString

	invalidCommandType = invalidRestrictedString

	invalidEventType = invalidRestrictedString
)

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gRPC Validator Middleware Test Suite")
}

func NewValidCreateCluster() *commanddata.CreateCluster {
	return &commanddata.CreateCluster{
		Name:             validName,
		ApiServerAddress: validApiServerAddress,
		CaCertBundle:     []byte(validCertificate),
	}
}

func NewValidClusterCreated() *eventdata.ClusterCreated {
	return &eventdata.ClusterCreated{
		Name:                validDisplayName,
		Label:               validName,
		ApiServerAddress:    validApiServerAddress,
		CaCertificateBundle: []byte(validCertificate),
	}
}

func NewValidClusterCreatedV2() *eventdata.ClusterCreatedV2 {
	return &eventdata.ClusterCreatedV2{
		Name:                validName,
		DisplayName:         validDisplayName,
		ApiServerAddress:    validApiServerAddress,
		CaCertificateBundle: []byte(validCertificate),
	}
}

func NewValidUpdateCluster() *commanddata.UpdateCluster {
	return &commanddata.UpdateCluster{
		Name:             &wrapperspb.StringValue{Value: validDisplayName},
		ApiServerAddress: &wrapperspb.StringValue{Value: validApiServerAddress},
		CaCertBundle:     []byte(validCertificate),
	}
}

func NewValidClusterUpdated() *eventdata.ClusterUpdated {
	return &eventdata.ClusterUpdated{
		DisplayName:         validDisplayName,
		ApiServerAddress:    validApiServerAddress,
		CaCertificateBundle: []byte(validCertificate),
	}
}

func NewValidCreateTenantClusterBindingCommandData() *commanddata.CreateTenantClusterBindingCommandData {
	return &commanddata.CreateTenantClusterBindingCommandData{
		TenantId:  validUUID,
		ClusterId: validUUID,
	}
}

func NewValidTenantClusterBindingCreated() *eventdata.TenantClusterBindingCreated {
	return &eventdata.TenantClusterBindingCreated{
		TenantId:  validUUID,
		ClusterId: validUUID,
	}
}

func NewValidCreateTenantCommandData() *commanddata.CreateTenantCommandData {
	return &commanddata.CreateTenantCommandData{
		Name:   validDisplayName,
		Prefix: validTenantPrefix,
	}
}

func NewValidTenantCreated() *eventdata.TenantCreated {
	return &eventdata.TenantCreated{
		Name:   validDisplayName,
		Prefix: validTenantPrefix,
	}
}

func NewValidUpdateTenantCommandData() *commanddata.UpdateTenantCommandData {
	return &commanddata.UpdateTenantCommandData{
		Name: &wrapperspb.StringValue{Value: validDisplayName},
	}
}

func NewValidTenantUpdated() *eventdata.TenantUpdated {
	return &eventdata.TenantUpdated{
		Name: &wrapperspb.StringValue{Value: validDisplayName},
	}
}

func NewValidCreateUserCommandData() *commanddata.CreateUserCommandData {
	return &commanddata.CreateUserCommandData{
		Name:  validDisplayName,
		Email: validEmail,
	}
}

func NewValidUserCreated() *eventdata.UserCreated {
	return &eventdata.UserCreated{
		Name:  validDisplayName,
		Email: validEmail,
	}
}

func NewValidCreateUserRoleBindingCommandData() *commanddata.CreateUserRoleBindingCommandData {
	return &commanddata.CreateUserRoleBindingCommandData{
		UserId:   validUUID,
		Role:     validRole,
		Scope:    validScope,
		Resource: wrapperspb.String(validUUID),
	}
}

func NewValidUserRoleAdded() *eventdata.UserRoleAdded {
	return &eventdata.UserRoleAdded{
		UserId:   validUUID,
		Role:     validRole,
		Scope:    validScope,
		Resource: validUUID,
	}
}

func NewValidPermissionModel() *domain.PermissionModel {
	return &domain.PermissionModel{
		Roles:  []string{validRole, validRole, validRole},
		Scopes: []string{validScope, validScope, validScope},
	}
}

func NewValidCommand() *commands.Command {
	return &commands.Command{
		Id:   validUUID,
		Type: validCommandType,
		Data: &anypb.Any{},
	}
}

func NewValidCommandReply() *eventsourcing.CommandReply {
	return &eventsourcing.CommandReply{
		AggregateId: validUUID,
	}
}

func NewValidEvent() *eventsourcing.Event {
	return &eventsourcing.Event{
		Type:          validEventType,
		AggregateId:   validUUID,
		AggregateType: validAggregateType,
	}
}

func NewValidEventFilter() *eventsourcing.EventFilter {
	return &eventsourcing.EventFilter{
		AggregateId:   &wrapperspb.StringValue{Value: validUUID},
		AggregateType: &wrapperspb.StringValue{Value: validAggregateType},
	}
}

func NewValidClusterAuthTokenRequest() *gateway.ClusterAuthTokenRequest {
	return &gateway.ClusterAuthTokenRequest{
		ClusterId: validUUID,
		Role:      validRole,
	}
}
