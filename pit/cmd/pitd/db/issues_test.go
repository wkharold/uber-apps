package db

import (
	"database/sql"
	"fmt"
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

var (
	issueone = Issue{
		id: 2001, description: "issue one", priority: 1, status: Open, project: 101, reporter: "fred@testrock.org",
	}
	issuetwo = Issue{
		id: 2002, description: "issue two", priority: 2, status: Open, project: 102, reporter: "barney@testrock.org",
	}
	issuethree = Issue{
		id: 2003, description: "issue three", priority: 2, status: Closed, project: 102, reporter: "barney@testrock.org",
	}
	issuefour = Issue{
		id: 2004, description: "issue four", priority: 3, status: Returned, project: 104, reporter: "fred@testrock.org",
	}
	issuefive = Issue{
		id: 2005, description: "issue five", priority: 2, status: Open, project: 103, reporter: "fred@testrock.org",
	}
	issuesix = Issue{
		id: 2006, description: "issue six", priority: 4, status: Closed, project: 106, reporter: "wilma@testrock.org",
	}
)

var (
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
)

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
	                      (2001, "issue one", 1, "OPEN", 101, 1007),
						  (2002, "issue two", 2, "OPEN", 102, 1008),
						  (2003, "issue three", 2, "CLOSED", 102, 1008),
						  (2004, "issue four", 3, "RETURNED", 104, 1007),
						  (2005, "issue five", 2, "OPEN", 103, 1007),
						  (2006, "issue six", 4, "CLOSED", 106, 1009);`); err != nil {
		panic(fmt.Sprintf("cannot setup issues table: [%+v]", err))
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
	                      (2001, "issue one", 1, "OPEN", 101, 1007);`); err != nil {
		panic(fmt.Sprintf("cannot setup issues table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}
