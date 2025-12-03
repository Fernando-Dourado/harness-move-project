package services

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TEST_VALID_PIPELINE_YAML = "pipeline:\n  name: AnyFilename 90\n  identifier: AnyFilename_90\n  projectIdentifier: FernandoD\n  orgIdentifier: default\n  tags: {}\n  stages:\n    - stage:\n        name: Stage 1\n        identifier: Stage_1\n        description: \"\"\n        type: Custom\n        spec:\n          execution:\n            steps:\n              - step:\n                  type: ShellScript\n                  name: Random of 15s\n                  identifier: Random_of_15s\n                  spec:\n                    shell: Bash\n                    onDelegate: true\n                    source:\n                      type: Inline\n                      spec:\n                        script: |-\n                          R=$(($RANDOM%15))\n                          echo $R\n                          sleep $R\n                    environmentVariables: []\n                    outputVariables: []\n                  timeout: 30m\n              - step:\n                  type: ShellScript\n                  name: SQL Injection\n                  identifier: SQL_Injection\n                  spec:\n                    shell: Bash\n                    onDelegate: true\n                    source:\n                      type: Inline\n                      spec:\n                        script: |-\n                          export SQL1=\"SELECT 1 FROM DUAL\"\n                          export SQL2=\"TABLE ACCOUNT\"\n\n                          echo $SQL1\n                          echo $SQL2\n                    environmentVariables: []\n                    outputVariables: []\n                  timeout: 10m\n        tags: {}\n"
const MISSING_ORG_YAML = "service:\n  name: pi-svc\n  identifier: pisvc\n  serviceDefinition:\n    type: Kubernetes\n    spec:\n      manifests:\n        - manifest:\n            identifier: job\n            type: K8sManifest\n            spec:\n              store:\n                type: Harness\n                spec:\n                  files:\n                    - /pi-svc/job.yaml\n              skipResourceVersioning: false\n              enableDeclarativeRollback: false\n  gitOpsEnabled: false\n  projectIdentifier: FernandoD\n"
const MISSING_PROJECT_YAML = "service:\n  name: pi-svc\n  identifier: pisvc\n  serviceDefinition:\n    type: Kubernetes\n    spec:\n      manifests:\n        - manifest:\n            identifier: job\n            type: K8sManifest\n            spec:\n              store:\n                type: Harness\n                spec:\n                  files:\n                    - /pi-svc/job.yaml\n              skipResourceVersioning: false\n              enableDeclarativeRollback: false\n  gitOpsEnabled: false\n  orgIdentifier: default\n"

func TestCreateYaml(t *testing.T) {
	yaml := createYaml(TEST_VALID_PIPELINE_YAML, "default", "FernandoD", "non_default", "DouradoF", "sourceAccountId", "targetAccountId")

	assert.True(t, strings.Contains(yaml, "orgIdentifier: non_default"), "The orgIdentifier not replaced")
	assert.True(t, strings.Contains(yaml, "projectIdentifier: DouradoF"), "The projectIdentifier not replaced")
}

func TestCreateYaml_MissingOrg(t *testing.T) {
	expectedOrg := "  orgIdentifier: non_default\n"

	yaml := createYaml(MISSING_ORG_YAML, "default", "FernandoD", "non_default", "DouradoF", "sourceAccountId", "targetAccountId")
	assert.True(t, strings.Contains(yaml, expectedOrg), "The orgIdentifier is missing")
}

func TestCreateYaml_MissingProject(t *testing.T) {
	expectedProject := "  projectIdentifier: DouradoF\n"

	yaml := createYaml(MISSING_PROJECT_YAML, "default", "FernandoD", "non_default", "DouradoF", "sourceAccountId", "targetAccountId")
	assert.True(t, strings.Contains(yaml, expectedProject), "The projectIdentifier is missing")
}

func TestRemoveNewLine(t *testing.T) {
	v1 := "There are no eligible delegates available in the account to execute the task.\n\n\n"
	v2 := removeNewLine(v1)
	assert.Equal(t, "There are no eligible delegates available in the account to execute the task.", v2)

	v3 := removeNewLine(v2)
	assert.Equal(t, v2, v3)
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"devEksCluster", "dev_eks_cluster"},
		{"prodEksCluster", "prod_eks_cluster"},
		{"uatEksCluster", "uat_eks_cluster"},
		{"backupPath", "backup_path"},
		{"windowsRuntimePath", "windows_runtime_path"},
		{"engineeringEdamame", "engineering_edamame"},
		{"simple", "simple"},
		{"ABC", "a_b_c"},
		{"already_snake", "already_snake"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransformIdentifiers(t *testing.T) {
	input := `infrastructureDefinition:
  name: dev-eks-cluster
  identifier: devEksCluster
  orgIdentifier: test-org
  projectIdentifier: new-project`

	expected := `infrastructureDefinition:
  name: dev-eks-cluster
  identifier: dev_eks_cluster
  orgIdentifier: test-org
  projectIdentifier: new-project`

	result := transformIdentifiers(input)
	assert.Equal(t, expected, result)
}

func TestTransformIdentifiers_MultipleIdentifiers(t *testing.T) {
	input := `artifacts:
  primary:
    primaryArtifactRef: engineeringEdamame
    sources:
      - identifier: engineeringEdamame
        type: Ecr`

	result := transformIdentifiers(input)
	assert.Contains(t, result, "identifier: engineering_edamame")
}

func TestMapEnvironment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"staging", "uat"},
		{"development", "dev"},
		{"dev", "dev"},
		{"devEks", "dev"},
		{"integration", "dev"},
		{"int", "dev"},
		{"quality", "qa"},
		{"qa", "qa"},
		{"production", "prod"},
		{"prod", "prod"},
		{"prod_pci", "prod_pci"},
		{"production_pci", "prod_pci"},
		{"uat", "uat"},
		{"uat_pci", "uat_pci"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapEnvironment(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapInfrastructure(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"dev", "dev_eks"},
		{"development", "dev_eks"},
		{"int", "dev_eks"},
		{"integration", "dev_eks"},
		{"qa", "dev_eks"},
		{"quality", "dev_eks"},
		{"staging", "uat_eks"},
		{"uat", "uat_eks"},
		{"uat_pci", "uat_pci_eks"},
		{"prod", "prod_eks"},
		{"production", "prod_eks"},
		{"prod_pci", "prod_pci_eks"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapInfrastructure(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransformEnvironmentRefs(t *testing.T) {
	input := `environment:
    environmentRef: staging
    infrastructureDefinitions:
      - identifier: uat_eks_cluster`

	result := transformEnvironmentRefs(input)
	assert.Contains(t, result, "environmentRef: uat")
}

func TestTransformInfrastructureRefs(t *testing.T) {
	input := `infrastructureRef: devEksCluster`
	result := transformInfrastructureRefs(input)
	assert.Contains(t, result, "infrastructureRef: dev_eks")
}
