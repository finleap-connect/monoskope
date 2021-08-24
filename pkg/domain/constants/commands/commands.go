package commands

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

const (
	// Event emitted when a User has been created
	CreateUser es.CommandType = "CreateUser"
	// Command to delete a User
	DeleteUser es.CommandType = "DeleteUser"

	// Command to create a new UserRoleBinding
	CreateUserRoleBinding es.CommandType = "CreateUserRoleBinding"
	// Command to delete a UserRoleBinding
	DeleteUserRoleBinding es.CommandType = "DeleteUserRoleBinding"

	// Event emitted when a Tenant has been created
	CreateTenant es.CommandType = "CreateTenant"
	// Command to update a Tenant
	UpdateTenant es.CommandType = "UpdateTenant"
	// Command to delete a Tenant
	DeleteTenant es.CommandType = "DeleteTenant"

	// Command to create a Cluster
	CreateCluster es.CommandType = "CreateCluster"
	// Command to delete a Cluster
	DeleteCluster es.CommandType = "DeleteCluster"
	// Command to update a Cluster
	UpdateCluster es.CommandType = "UpdateCluster"

	// Command to request a certificate
	RequestCertificate es.CommandType = "RequestCertificate"
)
