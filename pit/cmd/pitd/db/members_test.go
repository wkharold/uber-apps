package db

import (
	"database/sql"
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

type contributesToTest struct {
	description string
	fn          func(context.Context) ([]Project, error)
	ctxfn       func() context.Context
	expected    []Project
	err         error
}

type findAllMembersTest struct {
	description string
	fn          func(Members, context.Context) ([]Member, error)
	ctxfn       func() context.Context
	expected    []Member
	err         error
}

type findMembersByEmailTest struct {
	description string
	fn          func(Members, context.Context, string) (Member, error)
	email       string
	ctxfn       func() context.Context
	expected    Member
	err         error
}

type findMembersByIDTest struct {
	description string
	fn          func(Members, context.Context, int) (Member, error)
	id          int
	ctxfn       func() context.Context
	expected    Member
	err         error
}

type memberIssuesTest struct {
	description string
	fn          func(context.Context) ([]Issue, error)
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

var (
	bob   = Member{id: 1003, email: "bob@members.com"}
	carol = Member{id: 1004, email: "carol@members.com"}
	ted   = Member{id: 1005, email: "ted@members.com"}
	alice = Member{id: 1006, email: "alice@members.com"}
)

var (
	assignmentsTests = []memberIssuesTest{
		{"Assignments empty tables", bob.Assignments, emptytables, []Issue{}, nil},
		{"Assignments one issue", bob.Assignments, alltheissues, []Issue{issuetwo}, nil},
		{"Assignments", carol.Assignments, alltheissues, []Issue{issuefive, issuesix}, nil},
	}
	contributesToTests = []contributesToTest{
		{"ContributesTo no contributions", bob.ContributesTo, contributions, []Project{}, nil},
		{"ContributesTo one project", carol.ContributesTo, contributions, []Project{pone}, nil},
		{"ContributesTo many projects", alice.ContributesTo, contributions, []Project{pone, pthree}, nil},
	}
	findAllMemberTests = []findAllMembersTest{
		{"FindAll no members", Members.FindAll, emptytables, []Member{}, nil},
		{"FindAll one member", Members.FindAll, onemember, []Member{bob}, nil},
		{"FindAll members", Members.FindAll, manymembers, []Member{carol, ted, alice}, nil},
	}
	findMembersByEmailTests = []findMembersByEmailTest{
		{"FindByEmail empty tables", Members.FindByEmail, "bob@members.com", emptytables, Member{}, sql.ErrNoRows},
		{"FindByEmail many members no match", Members.FindByEmail, "fred.c.dobbs@members.com", manymembers, Member{}, sql.ErrNoRows},
		{"FindByEmail one member", Members.FindByEmail, "bob@members.com", onemember, bob, nil},
		{"FindByEmail members", Members.FindByEmail, "ted@members.com", manymembers, ted, nil},
	}
	findMembersByIDTests = []findMembersByIDTest{
		{"FindByID empty tables", Members.FindByID, 1003, emptytables, Member{}, sql.ErrNoRows},
		{"FindByID many members no match", Members.FindByID, 2001, manymembers, Member{}, sql.ErrNoRows},
		{"FindByID one member", Members.FindByID, 1003, onemember, bob, nil},
		{"FindByID members", Members.FindByID, 1005, manymembers, ted, nil},
	}
	watchingTests = []memberIssuesTest{
		{"Watching empty tables", bob.Watching, emptytables, []Issue{}, nil},
		{"Watching one issue", alice.Watching, alltheissues, []Issue{issuefive}, nil},
		{"Watching", bob.Watching, alltheissues, []Issue{issuetwo, issuethree, issuesix}, nil},
	}
)

func TestMemberAssignments(t *testing.T) {
	for _, nt := range assignmentsTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		is, err := nt.fn(ctx)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil && err == nt.err:
			break
		default:
			if !sameissues(is, nt.expected) {
				t.Errorf("%s: got %+v, expected %+v", nt.description, is, nt.expected)
				break
			}
		}

		dropdb(db)
	}
}

func TestMemberContributesTo(t *testing.T) {
	for _, nt := range contributesToTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		ps, err := nt.fn(ctx)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil && err == nt.err:
			break
		default:
			if !sameprojects(ps, nt.expected) {
				t.Errorf("%s: got %+v, expected %+v", nt.description, ps, nt.expected)
				break
			}
		}

		dropdb(db)
	}
}

func TestMemberWatching(t *testing.T) {
	for _, nt := range watchingTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		is, err := nt.fn(ctx)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil && err == nt.err:
			break
		default:
			if !sameissues(is, nt.expected) {
				t.Errorf("%s: got %+v, expected %+v", nt.description, is, nt.expected)
				break
			}
		}

		dropdb(db)
	}
}

func TestFindAllMembers(t *testing.T) {
	for _, nt := range findAllMemberTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		ms, err := nt.fn(struct{}{}, ctx)
		if err != nil {
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
			dropdb(db)
			continue
		}

		if !samemembers(nt.expected, ms) {
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, ms)
			dropdb(db)
			continue
		}

		dropdb(db)
	}
}

func TestFindMembersByEmail(t *testing.T) {
	for _, nt := range findMembersByEmailTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := nt.fn(struct{}{}, ctx, nt.email)
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

func TestFindMembersByID(t *testing.T) {
	for _, nt := range findMembersByIDTests {
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

func contributions() context.Context {
	db := createdb("contributions")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io"),
						  (1003, "bob@members.com"),
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO contributors VALUES
						  (101, 1006),
						  (101, 1004),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func onemember() context.Context {
	db := createdb("onemember")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES (1003, "bob@members.com");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func manymembers() context.Context {
	db := createdb("manymembers")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1004, "carol@members.com"),
						  (1005, "ted@members.com"),
						  (1006, "alice@members.com");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}
