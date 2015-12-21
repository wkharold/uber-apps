package db

import (
	"database/sql"
	"fmt"
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
	expected    Issue
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
	description string
	fn          func(Issues, context.Context, string) ([]Issue, error)
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
		{"FindAll empty tables", Issues.FindAll, emptytables, []Issue{}, nil},
		{"FindAll no issues", Issues.FindAll, noissues, []Issue{}, nil},
		{"FindAll one issue", Issues.FindAll, oneissue, []Issue{issueone}, nil},
		{"FindAll issues", Issues.FindAll, alltheissues, []Issue{issueone, issuetwo, issuethree, issuefour, issuefive, issuesix}, nil},
	}
	findIssuesByIDTests = []findIssuesByIDTest{
		{"FindByID empty tables", Issues.FindByID, 42, emptytables, Issue{}, sql.ErrNoRows},
		{"FindByID one issue", Issues.FindByID, 2001, oneissue, issueone, nil},
		{"FindByID one issue no match", Issues.FindByID, 1001, oneissue, Issue{}, sql.ErrNoRows},
		{"FindByID issues", Issues.FindByID, 2001, alltheissues, issueone, nil},
		{"FindByID issues no match", Issues.FindByID, 1001, alltheissues, Issue{}, sql.ErrNoRows},
	}
	findIssuesByPriorityTests = []findIssuesByPriorityTest{
		{"FindByPriority empty tables", Issues.FindByPriority, 1, emptytables, []Issue{}, nil},
		{"FindByPriority one issue", Issues.FindByPriority, 1, oneissue, []Issue{issueone}, nil},
		{"FindByPriority one issue no match", Issues.FindByPriority, 3, oneissue, []Issue{}, nil},
		{"FindByPriority issues no match", Issues.FindByPriority, 5, alltheissues, []Issue{}, nil},
		{"FindByPriority issues one match", Issues.FindByPriority, 1, alltheissues, []Issue{issueone}, nil},
		{"FindByPriority issues", Issues.FindByPriority, 2, alltheissues, []Issue{issuetwo, issuethree, issuefive}, nil},
	}
	findIssuesByProjectTests = []findIssuesByProjectTest{
		{"FindByProject empty tables", Issues.FindByProject, 112, emptytables, []Issue{}, nil},
		{"FindByProject one issue", Issues.FindByProject, 101, oneissue, []Issue{issueone}, nil},
		{"FindByProject one issue no match", Issues.FindByProject, 112, oneissue, []Issue{}, nil},
		{"FindByProject issues no match", Issues.FindByProject, 212, alltheissues, []Issue{}, nil},
		{"FindByProject issues one match", Issues.FindByProject, 101, alltheissues, []Issue{issueone}, nil},
		{"FindByProject issues", Issues.FindByProject, 102, alltheissues, []Issue{issuetwo, issuethree}, nil},
	}
	findIssuesByReporterTests = []findIssuesByReporterTest{
		{"FindByReporter empty tables", Issues.FindByReporter, "fred.c.dobbs@sierra.madre", emptytables, []Issue{}, nil},
		{"FindByReporter one issue", Issues.FindByReporter, "fred@testrock.org", oneissue, []Issue{issueone}, nil},
		{"FindByReporter one issue no match", Issues.FindByReporter, "betty@testrock.org", oneissue, []Issue{}, nil},
		{"FindByReporter issues no match", Issues.FindByReporter, "betty@testrock.org", alltheissues, []Issue{}, nil},
		{"FindByReporter issues one match", Issues.FindByReporter, "wilma@testrock.org", alltheissues, []Issue{issuesix}, nil},
		{"FindByReporter issues", Issues.FindByReporter, "fred@testrock.org", alltheissues, []Issue{issueone, issuefour, issuefive}, nil},
	}
	findIssuesByStatusTests = []findIssuesByStatusTest{
		{"FindByStatus empty tables", Issues.FindByStatus, Closed, emptytables, []Issue{}, nil},
		{"FindByStatus one issue", Issues.FindByStatus, Open, oneissue, []Issue{issueone}, nil},
		{"FindByStatus one issue on match", Issues.FindByStatus, Returned, oneissue, []Issue{}, nil},
		{"FindByStatus issues no match", Issues.FindByStatus, "UNKNOWN", alltheissues, []Issue{}, nil},
		{"FindByStatus issues one match", Issues.FindByStatus, Returned, alltheissues, []Issue{issuefour}, nil},
		{"FindByStatus issues", Issues.FindByStatus, Closed, alltheissues, []Issue{issuethree, issuesix}, nil},
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

		is, err := nt.fn(struct{}{}, ctx)
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

func TestFindIssuesByID(t *testing.T) {
	for _, nt := range findIssuesByIDTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := nt.fn(struct{}{}, ctx, nt.id)
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

		is, err := nt.fn(struct{}{}, ctx, nt.priority)
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

		is, err := nt.fn(struct{}{}, ctx, nt.project)
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

		is, err := nt.fn(struct{}{}, ctx, nt.reporter)
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

		is, err := nt.fn(struct{}{}, ctx, nt.status)
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

func alltheissues() context.Context {
	db := createdb("alltheissues")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002),
						  (104, "project four", "fourth test project", 1001),
						  (105, "project five", "fifth test project", 1002),
						  (106, "project six", "sixth test project", 1001);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com"),
						  (1007, "fred@testrock.org"),
						  (1008, "barney@testrock.org"),
						  (1009, "wilma@testrock.org");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO contributors VALUES
						  (102, 1006),
						  (103, 1004),
						  (103, 1005),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO issues VALUES
	                      (2001, "issueone", "issue one", 1, "OPEN", 101, 1007),
						  (2002, "issuetwo", "issue two", 2, "OPEN", 102, 1008),
						  (2003, "issuethree", "issue three", 2, "CLOSED", 102, 1008),
						  (2004, "issuefour", "issue four", 3, "RETURNED", 104, 1007),
						  (2005, "issuefive", "issue five", 2, "OPEN", 103, 1007),
						  (2006, "issuesix", "issue six", 4, "CLOSED", 106, 1009);`); err != nil {
		panic(fmt.Sprintf("cannot setup issues table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO assignments VALUES
	                      (1003, 2002),
						  (1004, 2005),
						  (1004, 2006),
						  (1005, 2006),
						  (1006, 2006);`); err != nil {
		panic(fmt.Sprintf("cannot setup assignments table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO watchers VALUES
	                      (1006, 2005),
	                      (1003, 2002),
						  (1003, 2003),
						  (1003, 2006),
	                      (1004, 2003),
						  (1005, 2003);`); err != nil {
		panic(fmt.Sprintf("cannot setup watchers table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func noissues() context.Context {
	db := createdb("noissues")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002),
						  (104, "project four", "fourth test project", 1001),
					      (105, "project five", "fifth test project", 1002),
						  (106, "project six", "sixth test project", 1001);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com"),
						  (1007, "fred@testrock.org"),
						  (1008, "barney@testrock.org");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO contributors VALUES
						  (102, 1006),
						  (103, 1004),
						  (103, 1005),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func oneissue() context.Context {
	db := createdb("alltheissues")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002),
						  (104, "project four", "fourth test project", 1001),
					      (105, "project five", "fifth test project", 1002),
						  (106, "project six", "sixth test project", 1001);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com"),
						  (1007, "fred@testrock.org"),
						  (1008, "barney@testrock.org");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO contributors VALUES
						  (102, 1006),
						  (103, 1004),
						  (103, 1005),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO issues VALUES
	                      (2001, "issueone", "issue one", 1, "OPEN", 101, 1007);`); err != nil {
		panic(fmt.Sprintf("cannot setup issues table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}
