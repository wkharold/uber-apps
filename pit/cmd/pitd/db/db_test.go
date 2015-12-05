package db

import (
	"testing"

	"golang.org/x/net/context"
)

type setfindertest struct {
	description string
	ffn         func(Projects, ctx context.Context) ([]Project, error)
	ctx         context.Context
	expected    []Project
}

type findertest struct {
	description string
	ffn         func(Projects, ctx context.Context) (Project, error)
	ctx         context.Context
	expected    Project
}

var (
	setfindertests = []setfindertest{}
	findertests    = []findertest{}
)

func TestSetFinders(t *testing.T) {
	for _, _ = range setfindertests {
	}
}

func TestFinders(t *testing.T) {
	for _, _ = range findertests {
	}
}
