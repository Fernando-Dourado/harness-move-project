package main

import (
	"testing"

	"github.com/Fernando-Dourado/harness-move-project/operation"
	"github.com/stretchr/testify/assert"
)

func TestArgumentRule_EmptyTargetProject(t *testing.T) {
	mv := operation.Move{
		Source: operation.NoName{
			Project: "ProjectA",
		},
	}
	applyArgumentRules(&mv)

	assert.Equal(t, mv.Source.Project, mv.Target.Project)
}

func TestArgumentRule_NonEmptyTargetProject(t *testing.T) {
	mv := operation.Move{
		Source: operation.NoName{
			Project: "ProjectA",
		},
		Target: operation.NoName{
			Project: "ProjectB",
		},
	}
	applyArgumentRules(&mv)

	assert.NotEqual(t, mv.Source.Project, mv.Target.Project)
	assert.Equal(t, mv.Source.Project, "ProjectA")
	assert.Equal(t, mv.Target.Project, "ProjectB")
}
