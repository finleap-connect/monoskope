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

package events

import (
	ef "github.com/finleap-connect/monoskope/pkg/audit/eventformatter"
)

// DetailsFormats
const (
	UserCreatedDetailsFormat            ef.DetailsFormat = "“%s“ created user “%s“"
	UserRoleAddedDetailsFormat          ef.DetailsFormat = "“%s“ assigned the role “%s“ for scope “%s“ to user “%s“"
	UserDeletedDetailsFormat            ef.DetailsFormat = "“%s“ deleted user “%s“"
	UserRoleBindingDeletedDetailsFormat ef.DetailsFormat = "“%s“ removed the role “%s“ for scope “%s“ from user “%s“"

	TenantCreatedDetailsFormat               ef.DetailsFormat = "“%s“ created tenant “%s“ with prefix “%s“"
	TenantUpdatedDetailsFormat               ef.DetailsFormat = "“%s“ updated the Tenant"
	TenantClusterBindingCreatedDetailsFormat ef.DetailsFormat = "“%s“ bounded tenant “%s“ to cluster “%s”"
	TenantDeletedDetailsFormat               ef.DetailsFormat = "“%s“ deleted tenant “%s“"
	TenantClusterBindingDeletedDetailsFormat ef.DetailsFormat = "“%s“ deleted the bound between cluster “%s“ and tenant “%s“"

	ClusterCreatedDetailsFormat               ef.DetailsFormat = "“%s“ created cluster “%s“"
	ClusterCreatedV2DetailsFormat             ef.DetailsFormat = ClusterCreatedDetailsFormat
	ClusterBootstrapTokenCreatedDetailsFormat ef.DetailsFormat = "“%s“ created a cluster bootstrap token"
	ClusterUpdatedDetailsFormat               ef.DetailsFormat = "“%s“ updated the cluster"
	ClusterDeletedDetailsFormat               ef.DetailsFormat = "“%s“ deleted cluster “%s“"

	RequestIssuedDetailsFormat            ef.DetailsFormat = "“%s“ issued a certificate request"
	CertificateRequestedDetailsFormat     ef.DetailsFormat = "“%s“ requested a certificate"
	CertificateIssuedDetailsFormat        ef.DetailsFormat = "“%s“ issued a certificate"
	CertificateIssuingFailedDetailsFormat ef.DetailsFormat = "certificate request issuing faild for “%s“"
)
