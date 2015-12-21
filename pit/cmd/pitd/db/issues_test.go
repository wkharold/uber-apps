package db

import (
	"database/sql"
	"testing"

	"golang.org/x/net/context"
)

type assignIssueTest struct {
	description string
	issue       Issue
	member      Member
	ctxfn       func() context.Context
	assigned    []Member
	err         error
}

type findAllIssuesTest struct {
	description string
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type findIssueByIDTest struct {
	description string
	id          int
	ctxfn       func() context.Context
	expected    Issue
	err         error
}

type findIssuesByPriorityTest struct {
	description string
	priority    int
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type findIssuesByProjectTest struct {
	description string
	project     int
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type findIssuesByReporterTest struct {
	description string
	reporter    string
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type findIssuesByStatusTest struct {
	description string
	status      string
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type issueTest struct {
	description string
	fn          func(context.Context) ([]Member, error)
	ctxfn       func() context.Context
	expected    []Member
	err         error
}

var (
	issueone = Issue{
		id: 2001, name: "issueone", description: "issue one", priority: 1, status: Open, project: 101, reporter: "fred@testrock.org",
	}
	issuetwo = Issue{
		id: 2002, name: "issuetwo", description: "issue two", priority: 2, status: Open, project: 102, reporter: "barney@testrock.org",
	}
	issuethree = Issue{
		id: 2003, name: "issuethree", description: "issue three", priority: 2, status: Closed, project: 102, reporter: "barney@testrock.org",
	}
	issuefour = Issue{
		id: 2004, name: "issuefour", description: "issue four", priority: 3, status: Returned, project: 104, reporter: "fred@testrock.org",
	}
	issuefive = Issue{
		id: 2005, name: "issuefive", description: "issue five", priority: 2, status: Open, project: 103, reporter: "fred@testrock.org",
	}
	issuesix = Issue{
		id: 2006, name: "issuesix", description: "issue six", priority: 4, status: Closed, project: 106, reporter: "wilma@testrock.org",
	}
	issueseven = Issue{
		id: 2007, name: "issueseven", description: "issue seven", priority: 2, status: Open, project: 102, reporter: "barney@testrock.org",
	}
)

var (
	assignIssueTests = []assignIssueTest{
		{"Assign non-existent member", issueone, Member{id: 9999, email: "fred.c.dobbs@sierramadre.gld"}, alltheissues, []Member{}, ErrNoSuchMember},
		{"Assign non-contributing member", issueone, bob, alltheissues, []Member{}, ErrNonContributingMember},
		{"Assign first assignee", issuethree, alice, alltheissues, []Member{alice}, nil},
		{"Assign additional assignee", issuetwo, alice, alltheissues, []Member{bob, alice}, nil},
	}
	findAllIssuesTests = []findAllIssuesTest{
		{"FindAll empty tables", emptytables, []Issue{}, nil},
		{"FindAll no issues", noissues, []Issue{}, nil},
		{"FindAll one issue", oneissue, []Issue{issueone}, nil},
		{"FindAll issues", alltheissues, []Issue{issueone, issuetwo, issuethree, issuefour, issuefive, issuesix}, nil},
	}
	findIssueByIDTests = []findIssueByIDTest{
		{"FindByID empty tables", 42, emptytables, Issue{}, sql.ErrNoRows},
		{"FindByID one issue", 2001, oneissue, issueone, nil},
		{"FindByID one issue no match", 1001, oneissue, Issue{}, sql.ErrNoRows},
		{"FindByID issues", 2001, alltheissues, issueone, nil},
		{"FindByID issues no match", 1001, alltheissues, Issue{}, sql.ErrNoRows},
	}
	findIssuesByPriorityTests = []findIssuesByPriorityTest{
		{"FindByPriority empty tables", 1, emptytables, []Issue{}, nil},
		{"FindByPriority one issue", 1, oneissue, []Issue{issueone}, nil},
		{"FindByPriority one issue no match", 3, oneissue, []Issue{}, nil},
		{"FindByPriority issues no match", 5, alltheissues, []Issue{}, nil},
		{"FindByPriority issues one match", 1, alltheissues, []Issue{issueone}, nil},
		{"FindByPriority issues", 2, alltheissues, []Issue{issuetwo, issuethree, issuefive}, nil},
	}
	findIssuesByProjectTests = []findIssuesByProjectTest{
		{"FindByProject empty tables", 112, emptytables, []Issue{}, nil},
		{"FindByProject one issue", 101, oneissue, []Issue{issueone}, nil},
		{"FindByProject one issue no match", 112, oneissue, []Issue{}, nil},
		{"FindByProject issues no match", 212, alltheissues, []Issue{}, nil},
		{"FindByProject issues one match", 101, alltheissues, []Issue{issueone}, nil},
		{"FindByProject issues", 102, alltheissues, []Issue{issuetwo, issuethree}, nil},
	}
	findIssuesByReporterTests = []findIssuesByReporterTest{
		{"FindByReporter empty tables", "fred.c.dobbs@sierra.madre", emptytables, []Issue{}, nil},
		{"FindByReporter one issue", "fred@testrock.org", oneissue, []Issue{issueone}, nil},
		{"FindByReporter one issue no match", "betty@testrock.org", oneissue, []Issue{}, nil},
		{"FindByReporter issues no match", "betty@testrock.org", alltheissues, []Issue{}, nil},
		{"FindByReporter issues one match", "wilma@testrock.org", alltheissues, []Issue{issuesix}, nil},
		{"FindByReporter issues", "fred@testrock.org", alltheissues, []Issue{issueone, issuefour, issuefive}, nil},
	}
	findIssuesByStatusTests = []findIssuesByStatusTest{
		{"FindByStatus empty tables", Closed, emptytables, []Issue{}, nil},
		{"FindByStatus one issue", Open, oneissue, []Issue{issueone}, nil},
		{"FindByStatus one issue on match", Returned, oneissue, []Issue{}, nil},
		{"FindByStatus issues no match", "UNKNOWN", alltheissues, []Issue{}, nil},
		{"FindByStatus issues one match", Returned, alltheissues, []Issue{issuefour}, nil},
		{"FindByStatus issues", Closed, alltheissues, []Issue{issuethree, issuesix}, nil},
	}
	issueAssignedTests = []issueTest{
		{"Assigned empty tables", issueone.Assigned, emptytables, []Member{}, nil},
		{"Assigned no assignement", issueone.Assigned, alltheissues, []Member{}, nil},
		{"Assigned one assignement", issuetwo.Assigned, alltheissues, []Member{bob}, nil},
		{"Assigned", issuesix.Assigned, alltheissues, []Member{carol, ted, alice}, nil},
	}
	issueWatchingTests = []issueTest{
		{"Watching empty tables", issueone.Watching, emptytables, []Member{}, nil},
		{"Watching no watchers", issuefour.Watching, alltheissues, []Member{}, nil},
		{"Watching one watcher", issuefive.Watching, alltheissues, []Member{alice}, nil},
		{"Watching", issuethree.Watching, alltheissues, []Member{bob, carol, ted}, nil},
	}
)

func TestAssignIssue(t *testing.T) {
	for _, nt := range assignIssueTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		err := nt.issue.Assign(ctx, nt.member)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil:
			break
		default:
			ms, err := nt.issue.Assigned(ctx)
			if err != nil {
				t.Errorf("%s: unable to verify assignments: [%+v]", nt.description, err)
				break
			}

			if !samemembers(nt.assigned, ms) {
				t.Error("%s: expected %+v, got %+v", nt.description, nt.assigned, ms)
				break
			}
		}

		dropdb(db)
	}
}
func TestFindAllIssues(t *testing.T) {
	for _, nt := range findAllIssuesTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		is, err := FindAllIssues(ctx)
		switch {
		case err != nil:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case !sameissues(nt.expected, is):
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, is)
		default:
			break
		}

		dropdb(db)
	}
}

func TestFindIssueByID(t *testing.T) {
	for _, nt := range findIssueByIDTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := FindIssueByID(ctx, nt.id)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil && err == nt.err:
			break
		default:
			if p != nt.expected {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, p)
				break
			}
		}

		dropdb(db)
	}
}

func TestFindIssuesByPriority(t *testing.T) {
	for _, nt := range findIssuesByPriorityTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		is, err := FindIssuesByPriority(ctx, nt.priority)
		switch {
		case err != nil:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case !sameissues(nt.expected, is):
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, is)
		default:
			break
		}

		dropdb(db)
	}
}

func TestFindIssuesByProject(t *testing.T) {
	for _, nt := range findIssuesByProjectTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		is, err := FindIssuesByProject(ctx, nt.project)
		switch {
		case err != nil:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case !sameissues(nt.expected, is):
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, is)
		default:
			break
		}

		dropdb(db)
	}
}

func TestFindIssuesByReporter(t *testing.T) {
	for _, nt := range findIssuesByReporterTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		is, err := FindIssuesByReporter(ctx, nt.reporter)
		switch {
		case err != nil:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case !sameissues(nt.expected, is):
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, is)
		default:
			break
		}

		dropdb(db)
	}
}

func TestFindIssuesByStatus(t *testing.T) {
	for _, nt := range findIssuesByStatusTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		is, err := FindIssuesByStatus(ctx, nt.status)
		switch {
		case err != nil:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case !sameissues(nt.expected, is):
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, is)
		default:
			break
		}

		dropdb(db)
	}
}

func TestIssueAssigned(t *testing.T) {
	for _, nt := range issueAssignedTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		ps, err := nt.fn(ctx)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil && err == nt.err:
			break
		default:
			if !samemembers(ps, nt.expected) {
				t.Errorf("%s: got %+v, expected %+v", nt.description, ps, nt.expected)
				break
			}
		}

		dropdb(db)
	}
}

func TestIssueWatching(t *testing.T) {
	for _, nt := range issueWatchingTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		ps, err := nt.fn(ctx)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil && err == nt.err:
			break
		default:
			if !samemembers(ps, nt.expected) {
				t.Errorf("%s: got %+v, expected %+v", nt.description, ps, nt.expected)
				break
			}
		}

		dropdb(db)
	}
}
