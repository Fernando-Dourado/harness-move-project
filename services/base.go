package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/Fernando-Dourado/harness-move-project/model"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

// BaseURL is the default Harness API URL used when no specific URL is provided
const BaseURL = "https://app.harness.io"

var (
	ErrEntityNotFound = errors.New("entity not found")
)

type SourceRequest struct {
	Client  *resty.Client
	Token   string
	Account string
	Url     string
}

type TargetRequest struct {
	Client  *resty.Client
	Token   string
	Account string
	Url     string
}

type SourceTarget struct {
	SourceOrg     string
	SourceProject string
	TargetOrg     string
	TargetProject string
	ExportDir     string
}

type Operation interface {
	Move() error
}

type OperationFactory interface {
	NewProjectOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) ProjectContext

	NewVariableOperation(sourceApi *SourceRequest, targetApi *TargetRequest, st *SourceTarget) VariableContext
}

// toSnakeCase converts a camelCase or PascalCase string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// transformIdentifiers converts all identifier values in YAML to snake_case
func transformIdentifiers(yaml string) string {
	// Match identifier: value patterns (handles both with and without quotes)
	re := regexp.MustCompile(`(identifier:\s*)([a-zA-Z][a-zA-Z0-9]*)`)
	return re.ReplaceAllStringFunc(yaml, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) == 3 {
			return parts[1] + toSnakeCase(parts[2])
		}
		return match
	})
}

// mapEnvironment maps source environment names to target environment names
// Rules:
//   - staging -> uat
//   - starts with or contains 'dev' -> dev
//   - starts with 'int' -> dev
//   - starts with 'quality' or 'qa' -> qa
//   - starts with 'prod' -> prod (preserves prod_pci)
func mapEnvironment(env string) string {
	envLower := strings.ToLower(env)

	// Exact or prefix matches first
	if envLower == "staging" {
		return "uat"
	}
	if strings.HasPrefix(envLower, "int") {
		return "dev"
	}
	if strings.HasPrefix(envLower, "quality") || strings.HasPrefix(envLower, "qa") {
		return "qa"
	}
	if strings.HasPrefix(envLower, "prod") {
		// Preserve prod_pci distinction
		if strings.Contains(envLower, "pci") {
			return "prod_pci"
		}
		return "prod"
	}
	if strings.HasPrefix(envLower, "uat") {
		// Preserve uat_pci distinction
		if strings.Contains(envLower, "pci") {
			return "uat_pci"
		}
		return "uat"
	}
	// Contains 'dev' anywhere
	if strings.Contains(envLower, "dev") {
		return "dev"
	}

	// Default: return as-is (snake_case will be applied separately)
	return env
}

// mapInfrastructure maps environment to its corresponding infrastructure
func mapInfrastructure(env string) string {
	envLower := strings.ToLower(env)

	switch {
	case envLower == "dev" || strings.HasPrefix(envLower, "dev"):
		return "dev_eks"
	case envLower == "int" || strings.HasPrefix(envLower, "int"):
		return "dev_eks"
	case envLower == "qa" || strings.HasPrefix(envLower, "quality"):
		return "dev_eks"
	case envLower == "uat_pci" || strings.Contains(envLower, "uat") && strings.Contains(envLower, "pci"):
		return "uat_pci_eks"
	case envLower == "uat" || envLower == "staging" || strings.HasPrefix(envLower, "uat"):
		return "uat_eks"
	case envLower == "prod_pci" || strings.Contains(envLower, "prod") && strings.Contains(envLower, "pci"):
		return "prod_pci_eks"
	case envLower == "prod" || strings.HasPrefix(envLower, "prod"):
		return "prod_eks"
	default:
		return env + "_eks"
	}
}

// transformEnvironmentRefs maps environmentRef values to target environments
func transformEnvironmentRefs(yaml string) string {
	re := regexp.MustCompile(`(environmentRef:\s*)([a-zA-Z][a-zA-Z0-9_]*)`)
	return re.ReplaceAllStringFunc(yaml, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) == 3 {
			return parts[1] + mapEnvironment(parts[2])
		}
		return match
	})
}

// transformInfrastructureRefs maps infrastructureRef and infrastructure identifiers
func transformInfrastructureRefs(yaml string) string {
	// Handle infrastructureRef fields
	re := regexp.MustCompile(`(infrastructureRef:\s*)([a-zA-Z][a-zA-Z0-9_]*)`)
	yaml = re.ReplaceAllStringFunc(yaml, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) == 3 {
			return parts[1] + mapInfrastructure(parts[2])
		}
		return match
	})
	return yaml
}

// transformInfrastructureDefinitionIdentifiers maps identifiers inside infrastructureDefinitions blocks
func transformInfrastructureDefinitionIdentifiers(yaml string) string {
	// Match identifier inside infrastructureDefinitions array items
	// Pattern: "- identifier: value" within infrastructureDefinitions context
	re := regexp.MustCompile(`(infrastructureDefinitions:\s*\n(?:.*\n)*?\s*-\s*identifier:\s*)([a-zA-Z][a-zA-Z0-9_]*)`)
	yaml = re.ReplaceAllStringFunc(yaml, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) == 3 {
			return parts[1] + mapInfrastructure(parts[2])
		}
		return match
	})

	// Also handle single-line format: infrastructureDefinitions: - identifier: xxx
	re2 := regexp.MustCompile(`(infrastructureDefinitions:.*identifier:\s*)([a-zA-Z][a-zA-Z0-9_]*)`)
	yaml = re2.ReplaceAllStringFunc(yaml, func(match string) string {
		parts := re2.FindStringSubmatch(match)
		if len(parts) == 3 {
			return parts[1] + mapInfrastructure(parts[2])
		}
		return match
	})

	return yaml
}

func createYaml(yaml, sourceOrg, sourceProject, targetOrg, targetProject, sourceAccount, targetAccount string) string {
	var out string

	// Replace accountId if present
	if strings.Contains(yaml, "accountId: ") {
		out = strings.ReplaceAll(yaml, "accountId: "+sourceAccount, "accountId: "+targetAccount)
	} else {
		out = yaml
	}

	// Replace orgIdentifier
	if strings.Contains(out, "orgIdentifier: ") {
		out = strings.ReplaceAll(out, "orgIdentifier: "+sourceOrg, "orgIdentifier: "+targetOrg)
	} else {
		if len(targetOrg) > 0 {
			out = fmt.Sprintln(out, " orgIdentifier:", targetOrg)
		}
	}

	// Replace projectIdentifier
	if strings.Contains(out, "projectIdentifier: ") {
		out = strings.ReplaceAll(out, "projectIdentifier: "+sourceProject, "projectIdentifier: "+targetProject)
	} else {
		if len(targetProject) > 0 {
			out = fmt.Sprintln(out, " projectIdentifier:", targetProject)
		}
	}

	// Replace account. prefix with org. for connector refs and other account-scoped resources
	out = strings.ReplaceAll(out, "account.", "org.")

	// Map environment references to target environments
	out = transformEnvironmentRefs(out)

	// Map infrastructure references to target infrastructures
	out = transformInfrastructureRefs(out)

	// Transform infrastructureDefinitions identifiers (special case - uses env mapping for infra)
	out = transformInfrastructureDefinitionIdentifiers(out)

	// Transform all identifier values to snake_case
	out = transformIdentifiers(out)

	return out
}

func handleErrorResponse(resp *resty.Response) error {
	result := model.ErrorResponse{}
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}
	if result.Code == "ENTITY_NOT_FOUND" {
		return ErrEntityNotFound
	}
	if result.Code == "DUPLICATE_FIELD" {
		return nil
	}
	if strings.Contains(result.Message, "already exists") {
		return nil
	}
	return fmt.Errorf("%s: %s", result.Code, removeNewLine(result.Message))
}

func removeNewLine(value string) string {
	return strings.ReplaceAll(value, "\n", "")
}

func reportFailed(failed []string, description string) {
	if len(failed) > 0 {
		fmt.Println(color.RedString(fmt.Sprintf("Failed %s %d", description, len(failed))))
		fmt.Println(color.RedString(strings.Join(failed, "\n")))
	}
}

func createQueryParams(account, org, project string) map[string]string {
	params := map[string]string{
		"accountIdentifier": account,
	}
	if len(org) > 0 {
		params["orgIdentifier"] = org
	}
	if len(project) > 0 {
		params["projectIdentifier"] = project
	}
	return params
}

// saveToYAML saves data to a YAML file in the export directory
func saveToYAML(exportDir, resourceType, identifier string, data interface{}) error {
	if exportDir == "" {
		return nil // Export mode not enabled
	}

	// Create directory structure: exportDir/resourceType/
	dir := filepath.Join(exportDir, resourceType)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Convert identifier to snake_case for filename
	snakeCaseId := toSnakeCase(identifier)
	filename := strings.ReplaceAll(snakeCaseId, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filePath := filepath.Join(dir, filename+".yaml")

	// Convert to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	// Transform identifier values in YAML to snake_case
	yamlStr := transformIdentifiers(string(yamlData))

	// Write to file
	if err := os.WriteFile(filePath, []byte(yamlStr), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

// saveStringToYAML saves a YAML string directly to file
func saveStringToYAML(exportDir, resourceType, identifier, yamlContent string) error {
	if exportDir == "" {
		return nil // Export mode not enabled
	}

	// Create directory structure: exportDir/resourceType/
	dir := filepath.Join(exportDir, resourceType)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Convert identifier to snake_case for filename
	snakeCaseId := toSnakeCase(identifier)
	filename := strings.ReplaceAll(snakeCaseId, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filePath := filepath.Join(dir, filename+".yaml")

	// Transform identifier values in YAML content to snake_case
	transformedContent := transformIdentifiers(yamlContent)

	// Write to file
	if err := os.WriteFile(filePath, []byte(transformedContent), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

// isExportMode checks if export mode is enabled
func isExportMode(exportDir string) bool {
	return exportDir != ""
}
