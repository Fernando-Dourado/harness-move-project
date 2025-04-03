package main

import (
	"testing"

	"github.com/Fernando-Dourado/harness-move-project/operation"
	"github.com/stretchr/testify/assert"
)

func TestArgumentRule_EmptyTargetProject(t *testing.T) {
	mv := operation.Move{
		Source: operation.CopyConfig{
			Project: "ProjectA",
		},
	}
	applyArgumentRules(&mv)

	assert.Equal(t, mv.Source.Project, mv.Target.Project)
}

func TestArgumentRule_NonEmptyTargetProject(t *testing.T) {
	mv := operation.Move{
		Source: operation.CopyConfig{
			Project: "ProjectA",
		},
		Target: operation.CopyConfig{
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
		Source: operation.CopyConfig{
			Account: "AccountA",
			Token:   "TokenA",
		},
	}
	applyArgumentRules(&mv)

	assert.Equal(t, mv.Source.Project, mv.Target.Project)
	assert.Equal(t, mv.Source.Token, mv.Target.Token)
}
