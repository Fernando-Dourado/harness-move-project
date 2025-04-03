package main

import (
	"testing"

	"github.com/Fernando-Dourado/harness-move-project/operation"
	"github.com/stretchr/testify/assert"
)

func TestArgumentRule_EmptyTargetProject(t *testing.T) {
	mv := operation.Move{
		Source: operation.Config{
			Project: "ProjectA",
		},
	}
	applyArgumentRules(&mv)

	assert.Equal(t, mv.Source.Project, mv.Target.Project)
}

func TestArgumentRule_NonEmptyTargetProject(t *testing.T) {
	mv := operation.Move{
		Source: operation.Config{
			Project: "ProjectA",
		},
		Target: operation.Config{
			Project: "ProjectB",
		},
	}
	applyArgumentRules(&mv)

	assert.NotEqual(t, mv.Source.Project, mv.Target.Project)
	assert.Equal(t, mv.Source.Project, "ProjectA")
	assert.Equal(t, mv.Target.Project, "ProjectB")
}

func TestArgumentRule_EmptyTargetAccount(t *testing.T) {
	mv := operation.Move{
		Source: operation.Config{
			Account: "AccountA",
			Token:   "TokenA",
		},
	}
	applyArgumentRules(&mv)

	assert.Equal(t, mv.Source.Project, mv.Target.Project)
	assert.Equal(t, mv.Source.Token, mv.Target.Token)
}
