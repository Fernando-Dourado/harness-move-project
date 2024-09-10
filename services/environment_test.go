package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const TEST_ENV_ORG_QUOTES = "environment:\n  orgIdentifier: \"default\"\n  projectIdentifier: \"Ansible\"\n  identifier: \"eu\"\n  tags: {}\n  name: \"eu\"\n  type: \"Production\"\n"

func TestSanitizeEnvYaml_EnvOrgQuoted(t *testing.T) {
	yaml := sanitizeEnvYaml(TEST_ENV_ORG_QUOTES)
	assert.NotContains(t, yaml, "\"")
}
