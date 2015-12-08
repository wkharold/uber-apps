package db

import (
	"testing"

	"golang.org/x/net/context"
)

type findAllIssuesTest struct {
	description string
	fn          func(Issues, context.Context) ([]Issue, error)
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type findIssuesByIDTest struct {
	description string
	fn          func(Issues, context.Context, int) (Issue, error)
	id          int
	ctxfn       func() context.Context
	expected    Issues
	err         error
}

type findIssuesByPriorityTest struct {
	description string
	fn          func(Issues, context.Context, int) ([]Issue, error)
	priority    int
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type findIssuesByProjectTest struct {
	description string
	fn          func(Issues, context.Context, int) ([]Issue, error)
	project     int
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type findIssuesByReporterTest struct {
	description string
	fn          func(Issues, context.Context, string) ([]Issue, error)
	reporter    string
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type findIssuesByStatusTest struct {
}

var (
	findAllIssuesTests        = []findAllIssuesTest{}
	findIssuesByIDTests       = []findIssuesByIDTest{}
	findIssuesByPriorityTests = []findIssuesByPriorityTest{}
	findIssuesByProjectTests  = []findIssuesByProjectTest{}
	findIssuesByReporterTests = []findIssuesByReporterTest{}
	findIssuesByStatusTests   = []findIssuesByStatusTest{}
)

func TestFindAllIssues(t *testing.T) {
	t.Fail("Unimplemented")
}

func TestFindIssuesByID(t *testing.T) {
	t.Fail("Unimplemented")
}

func TestFindIssuesByPriority(t *testing.T) {
	t.Fail("Unimplemented")
}

func TestFindIssuesByProject(t *testing.T) {
	t.Fail("Unimplemented")
}

func TestFindIssuesByReport(t *testing.T) {
	t.Fail("Unimplemented")
}

func TestFindIssuesByStatus(t *testing.T) {
	t.Fail("Unimplemented")
}
