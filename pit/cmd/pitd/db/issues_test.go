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
)

var (
	findAllIssuesTests = []findAllIssuesTest{
		{"FindAll empty tables", Issues.FindAll, emptytables, []Issue{}, nil},
		{"FindAll no issues", Issues.FindAll, noissues, []Issue{}, nil},
		{"FindAll one issue", Issues.FindAll, oneissue, []Issue{issueone}, nil},
		{"FindAll issues", Issues.FindAll, alltheissues, []Issue{issueone}, nil},
	}
	findIssuesByIDTests = []findIssuesByIDTest{
		{"FindByID empty tables", Issues.FindByID, 42, emptytables, Issue{}, sql.ErrNoRows},
		{"FindByID one issue", Issues.FindByID, 2001, oneissue, issueone, nil},
	}
	findIssuesByPriorityTests = []findIssuesByPriorityTest{
		{"FindByPriority empty tables", Issues.FindByPriority, 1, emptytables, []Issue{}, nil},
		{"FindByPriority one issue", Issues.FindByPriority, 1, oneissue, []Issue{issueone}, nil},
		{"FindByPriority one issue no match", Issues.FindByPriority, 3, oneissue, []Issue{}, nil},
	}
	findIssuesByProjectTests = []findIssuesByProjectTest{
		{"FindByProject empty tables", Issues.FindByProject, 112, emptytables, []Issue{}, nil},
		{"FindByProject one issue", Issues.FindByProject, 101, oneissue, []Issue{issueone}, nil},
		{"FindByProject one issue no match", Issues.FindByProject, 112, oneissue, []Issue{}, nil},
	}
	findIssuesByReporterTests = []findIssuesByReporterTest{
		{"FindByReporter empty tables", Issues.FindByReporter, "fred.c.dobbs@sierra.madre", emptytables, []Issue{}, nil},
		{"FindByReporter one issue", Issues.FindByReporter, "fred@testrock.org", oneissue, []Issue{issueone}, nil},
		{"FindByReporter one issue no match", Issues.FindByReporter, "betty@testrock.org", oneissue, []Issue{}, nil},
	}
	findIssuesByStatusTests = []findIssuesByStatusTest{
		{"FindByStatus empty tables", Issues.FindByStatus, Closed, emptytables, []Issue{}, nil},
		{"FindByStatus one issue", Issues.FindByStatus, Open, oneissue, []Issue{issueone}, nil},
		{"FindByStatus one issue on match", Issues.FindByStatus, Returned, oneissue, []Issue{}, nil},
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
	t.Fail()
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
