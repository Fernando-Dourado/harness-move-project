package operation

import (
	"errors"

	"github.com/Fernando-Dourado/harness-move-project/services"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestCreateProjectWhenNotRequired_AndEntityNotFound(t *testing.T) {

	sourceApi := &services.SourceRequest{}
	targetApi := &services.TargetRequest{}

	move := Move{
		Config: OperationConfig{
			CreateProject: false,
		},
	}

	expectedError := services.ErrEntityNotFound
	actualError := move.createProjectWhenRequired(sourceApi, targetApi, expectedError)

	assert.Equal(t, expectedError, actualError)
}

func TestCreateProjectWhenNotRequired_AndOtherError(t *testing.T) {

	sourceApi := &services.SourceRequest{}
	targetApi := &services.TargetRequest{}

	move := Move{
		Config: OperationConfig{
			CreateProject: false,
		},
	}

	expectedError := errors.New("some other error")
	actualError := move.createProjectWhenRequired(sourceApi, targetApi, expectedError)

	assert.Equal(t, expectedError, actualError)
}
