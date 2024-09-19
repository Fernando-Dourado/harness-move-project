package services

// Define counter variables
var apiCalls int = 0

var projects int = 0

var connectorsTotal int = 0

var connectorsMoved int = 0

var environmentsTotal int = 0

var environmentsMoved int = 0

var environmentGroupsTotal int = 0

var environmentGroupsMoved int = 0

var featureFlagsTotal int = 0

var featureFlagsMoved int = 0

var fileStoresTotal int = 0

var fileStoresMoved int = 0

var infrastructureTotal int = 0

var infrastructureMoved int = 0

var inputSetTotal int = 0

var inputSetMoved int = 0

var overridesTotal int = 0

var overridesMoved int = 0

var pipelinesTotal int = 0

var pipelinesMoved int = 0

var resourceGroupsTotal int = 0

var resourceGroupsMoved int = 0

var roleAssignmentsTotal int = 0

var roleAssignmentsMoved int = 0

var rolesTotal int = 0

var rolesMoved int = 0

var servicesTotal int = 0

var servicesMoved int = 0

var serviceAccountsTotal int = 0

var serviceAccountsMoved int = 0

var tagsTotal int = 0

var tagsMoved int = 0

var targetGroupsTotal int = 0

var targetGroupsMoved int = 0

var targetsTotal int = 0

var targetsMoved int = 0

var templatesTotal int = 0

var templatesMoved int = 0

var triggersTotal int = 0

var triggersMoved int = 0

var userGroupsTotal int = 0

var userGroupsMoved int = 0

var usersTotal int = 0

var usersMoved int = 0

var variablesTotal int = 0

var variablesMoved int = 0

// Define counter functions to increment and get values

// API Calls
func IncrementApiCalls() {
	apiCalls++
}

func GetApiCalls() int {
	return apiCalls
}

// Projects
func IncrementProjects() {
	projects++
}

func GetProjects() int {
	return projects
}

// Connectors
func IncrementConnectorsTotal() {
	connectorsTotal++
}

func GetConnectorsTotal() int {
	return connectorsTotal
}

func GetConnectorsMoved() int {
	return connectorsMoved
}
func IncrementConnectorsMoved() {
	connectorsMoved++
}

// Environments
func IncrementEnvironmentsTotal() {
	environmentsTotal++
}

func GetEnvironmentsTotal() int {
	return environmentsTotal
}

func IncrementEnvironmentsMoved() {
	environmentsMoved++
}

func GetEnvironmentsMoved() int {
	return environmentsMoved
}

// Environment Groups
func IncrementEnvironmentGroupsTotal() {
	environmentGroupsTotal++
}

func GetEnvironmentGroupsTotal() int {
	return environmentGroupsTotal
}

func IncrementEnvironmentGroupsMoved() {
	environmentGroupsMoved++
}

func GetEnvironmentGroupsMoved() int {
	return environmentGroupsMoved
}

// Feature Flags
func IncrementFeatureFlagsTotal() {
	featureFlagsTotal++
}

func GetFeatureFlagsTotal() int {
	return featureFlagsTotal
}

func IncrementFeatureFlagsMoved() {
	featureFlagsMoved++
}

func GetFeatureFlagsMoved() int {
	return featureFlagsMoved
}

// File Stores
func IncrementFileStoresTotal() {
	fileStoresTotal++
}

func GetFileStoresTotal() int {
	return fileStoresTotal
}

func IncrementFileStoresMoved() {
	fileStoresMoved++
}

func GetFileStoresMoved() int {
	return fileStoresMoved
}

// Input Sets
func IncrementInputSetsTotal() {
	inputSetTotal++
}

func GetInputSetsTotal() int {
	return inputSetTotal
}

func IncrementInputSetsMoved() {
	inputSetMoved++
}

func GetInputSetsMoved() int {
	return inputSetMoved
}

// Infrastructure
func IncrementInfrastructureTotal() {
	infrastructureTotal++
}

func GetInfrastructureTotal() int {
	return infrastructureTotal
}

func IncrementInfrastructureMoved() {
	infrastructureMoved++
}

func GetInfrastructureMoved() int {
	return infrastructureMoved
}

// Overrides
func IncrementOverridesTotal() {
	overridesTotal++
}

func GetOverridesTotal() int {
	return overridesTotal
}

func IncrementOverridesMoved() {
	overridesMoved++
}

func GetOverridesMoved() int {
	return overridesMoved
}

// Pipelines
func IncrementPipelinesTotal() {
	pipelinesTotal++
}

func GetPipelinesTotal() int {
	return pipelinesTotal
}

func IncrementPipelinesMoved() {
	pipelinesMoved++
}

func GetPipelinesMoved() int {
	return pipelinesMoved
}

// Resource Groups
func IncrementResourceGroupsTotal() {
	resourceGroupsTotal++
}

func GetResourceGroupsTotal() int {
	return resourceGroupsTotal
}

func IncrementResourceGroupsMoved() {
	resourceGroupsMoved++
}

func GetResourceGroupsMoved() int {
	return resourceGroupsMoved
}

// Role Assignments
func IncrementRoleAssignmentsTotal() {
	roleAssignmentsTotal++
}

func GetRoleAssignmentsTotal() int {
	return roleAssignmentsTotal
}

func IncrementRoleAssignmentsMoved() {
	roleAssignmentsMoved++
}

func GetRoleAssignmentsMoved() int {
	return roleAssignmentsMoved
}

// Roles
func IncrementRolesTotal() {
	rolesTotal++
}

func GetRolesTotal() int {
	return rolesTotal
}

func IncrementRolesMoved() {
	rolesMoved++
}

func GetRolesMoved() int {
	return rolesMoved
}

// Services
func IncrementServicesTotal() {
	servicesTotal++
}

func GetServicesTotal() int {
	return servicesTotal
}

func IncrementServicesMoved() {
	servicesMoved++
}

func GetServicesMoved() int {
	return servicesMoved
}

// Service Accounts
func IncrementServiceAccountsTotal() {
	serviceAccountsTotal++
}

func GetServiceAccountsTotal() int {
	return serviceAccountsTotal
}

func IncrementServiceAccountsMoved() {
	serviceAccountsMoved++
}

func GetServiceAccountsMoved() int {
	return serviceAccountsMoved
}

// Tags
func IncrementTagsTotal() {
	tagsTotal++
}

func GetTagsTotal() int {
	return tagsTotal
}

func IncrementTagsMoved() {
	tagsMoved++
}

func GetTagsMoved() int {
	return tagsMoved
}

// Target Groups
func IncrementTargetGroupsTotal() {
	targetGroupsTotal++
}

func GetTargetGroupsTotal() int {
	return targetGroupsTotal
}

func IncrementTargetGroupsMoved() {
	targetGroupsMoved++
}

func GetTargetGroupsMoved() int {
	return targetGroupsMoved
}

// Targets
func IncrementTargetsTotal() {
	targetsTotal++
}

func GetTargetsTotal() int {
	return targetsTotal
}

func IncrementTargetsMoved() {
	targetsMoved++
}

func GetTargetsMoved() int {
	return targetsMoved
}

// Templates
func IncrementTemplatesTotal() {
	templatesTotal++
}

func GetTemplatesTotal() int {
	return templatesTotal
}

func IncrementTemplatesMoved() {
	templatesMoved++
}

func GetTemplatesMoved() int {
	return templatesMoved
}

// Triggers
func IncrementTriggersTotal() {
	triggersTotal++
}

func GetTriggersTotal() int {
	return triggersTotal
}

func IncrementTriggersMoved() {
	triggersMoved++
}

func GetTriggersMoved() int {
	return triggersMoved
}

// User Groups
func IncrementUserGroupsTotal() {
	userGroupsTotal++
}

func GetUserGroupsTotal() int {
	return userGroupsTotal
}

func IncrementUserGroupsMoved() {
	userGroupsMoved++
}

func GetUserGroupsMoved() int {
	return userGroupsMoved
}

// Users
func IncrementUsersTotal() {
	usersTotal++
}

func GetUsersTotal() int {
	return usersTotal
}

func IncrementUsersMoved() {
	usersMoved++
}

func GetUsersMoved() int {
	return usersMoved
}

// Variables
func IncrementVariablesTotal() {
	variablesTotal++
}

func GetVariablesTotal() int {
	return variablesTotal
}

func IncrementVariablesMoved() {
	variablesMoved++
}

func GetVariablesMoved() int {
	return variablesMoved
}
