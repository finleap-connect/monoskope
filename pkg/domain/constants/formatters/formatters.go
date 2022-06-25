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

package formatters

import (
	"fmt"
	"time"
)

type DetailsFormat string

const (
	TimeFormat       = time.RFC822
	LeftQuoteSymbol  = "“"
	RightQuoteSymbol = "“"

	UserCreatedDetailsFormat            DetailsFormat = "“%s“ created user “%s“"
	UserUpdatedDetailsFormat            DetailsFormat = "“%s“ updated the User"
	UserRoleAddedDetailsFormat          DetailsFormat = "“%s“ assigned the role “%s“ for scope “%s“ to user “%s“"
	UserDeletedDetailsFormat            DetailsFormat = "“%s“ deleted user “%s“"
	UserRoleBindingDeletedDetailsFormat DetailsFormat = "“%s“ removed the role “%s“ for scope “%s“ from user “%s“"

	TenantCreatedDetailsFormat               DetailsFormat = "“%s“ created tenant “%s“ with prefix “%s“"
	TenantUpdatedDetailsFormat               DetailsFormat = "“%s“ updated the Tenant"
	TenantClusterBindingCreatedDetailsFormat DetailsFormat = "“%s“ bounded tenant “%s“ to cluster “%s”"
	TenantDeletedDetailsFormat               DetailsFormat = "“%s“ deleted tenant “%s“"
	TenantClusterBindingDeletedDetailsFormat DetailsFormat = "“%s“ deleted the bound between cluster “%s“ and tenant “%s“"

	ClusterCreatedDetailsFormat   DetailsFormat = "“%s“ created cluster “%s“"
	ClusterCreatedV2DetailsFormat DetailsFormat = ClusterCreatedDetailsFormat
	ClusterUpdatedDetailsFormat   DetailsFormat = "“%s“ updated the cluster"
	ClusterDeletedDetailsFormat   DetailsFormat = "“%s“ deleted cluster “%s“"

	RequestIssuedDetailsFormat            DetailsFormat = "“%s“ issued a certificate request"
	CertificateRequestedDetailsFormat     DetailsFormat = "“%s“ requested a certificate"
	CertificateIssuedDetailsFormat        DetailsFormat = "“%s“ issued a certificate"
	CertificateIssuingFailedDetailsFormat DetailsFormat = "certificate request issuing faild for “%s“"

	UserCreatedOverviewDetailsFormat            DetailsFormat = "“%s“ was created by “%s“ at “%s“"
	UserDeletedOverviewDetailsFormat            DetailsFormat = " and was deleted by “%s“ at “%s“"
	UserRoleBindingOverviewDetailsFormat        DetailsFormat = "- %s %s\n"
	TenantUserRoleBindingOverviewDetailsFormat  DetailsFormat = "- %s (%s)\n"
	ClusterUserRoleBindingOverviewDetailsFormat DetailsFormat = TenantUserRoleBindingOverviewDetailsFormat
)

// Sprint returns the resulting string after formatting.
func (f DetailsFormat) Sprint(args ...interface{}) string {
	return fmt.Sprintf(string(f), args...)
}

func Quote(str string) string {
	return LeftQuoteSymbol + str + RightQuoteSymbol
}
