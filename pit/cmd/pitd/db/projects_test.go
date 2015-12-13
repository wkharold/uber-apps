package db

import (
	"database/sql"
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

type contributorstest struct {
	description string
	fn          func(context.Context) ([]Member, error)
	ctxfn       func() context.Context
	expected    []Member
	err         error
}

type findAllProjectsTest struct {
	description string
	fn          func(Projects, context.Context) ([]Project, error)
	ctxfn       func() context.Context
	expected    []Project
	err         error
}

type findProjectsByOwnerTest struct {
	description string
	fn          func(Projects, context.Context, string) ([]Project, error)
	owner       string
	ctxfn       func() context.Context
	expected    []Project
	err         error
}

type findProjectsByIDTest struct {
	description string
	fn          func(Projects, context.Context, int) (Project, error)
	id          int
	ctxfn       func() context.Context
	expected    Project
	err         error
}

type findProjectsByNameTest struct {
	description string
	fn          func(Projects, context.Context, string) (Project, error)
	name        string
	ctxfn       func() context.Context
	expected    Project
	err         error
}

type newProjectTest struct {
	description string
	fn          func(context.Context, string, string, string) (Project, error)
	name        string
	desc        string
	owner       string
	id          int
	ctxfn       func() context.Context
	expected    Project
	collection  []Project
	err         error
}

type openProjectIssueTest struct {
	description string
	fn          func(context.Context, string, string, string, int) (Issue, error)
	name        string
	desc        string
	reporter    string
	priority    int
	id          int
	ctxfn       func() context.Context
	expected    Issue
	collection  []Issue
	err         error
}

var (
	issueonep2 = Issue{
		id: 2012, name: "issueone", description: "issue one", priority: 2, status: Open, project: 102, reporter: "fred@testrock.org",
	}
	issueonep3 = Issue{
		id: 2011, name: "issueone", description: "issue one", priority: 2, status: Open, project: 103, reporter: "barne@testrock.org",
	}
	pone = Project{
		id: 101, name: "project one", description: "first test project", owner: "owner@test.net",
	}
	ptwo = Project{
		id: 102, name: "project two", description: "second test project", owner: "owner@test.net",
	}
	pthree = Project{
		id: 103, name: "project three", description: "third test project", owner: "owner@test.io",
	}
	pfour = Project{
		id: 104, name: "project four", description: "fourth test project", owner: "owner@test.net",
	}
)

var (
	contributorstests = []contributorstest{
		{"Contributors no contributors", pone.Contributors, contributors, []Member{}, nil},
		{"Contributors single contributor", ptwo.Contributors, contributors, []Member{alice}, nil},
		{"Contributors multiple", pthree.Contributors, contributors, []Member{carol, ted, alice}, nil},
	}
	findAllProjectsTests = []findAllProjectsTest{
		{"FindAll from empty tables", Projects.FindAll, emptytables, []Project{}, nil},
		{"FindAll one project", Projects.FindAll, oneproject, []Project{pone}, nil},
		{"FindAll multiple projects", Projects.FindAll, manyprojects, []Project{pone, ptwo, pthree}, nil},
	}
	findProjectsByOwnerTests = []findProjectsByOwnerTest{
		{"FindByOwner from empty tables", Projects.FindByOwner, "owner@test.net", emptytables, []Project{}, nil},
		{"FindByOwner one project no match", Projects.FindByOwner, "owner@test.org", oneproject, []Project{}, nil},
		{"FindByOwner multiple projects no match", Projects.FindByOwner, "owner@test.com", manyprojects, []Project{}, nil},
		{"FindByOwner one project", Projects.FindByOwner, "owner@test.net", oneproject, []Project{pone}, nil},
		{"FindByOwner multiple projects one match", Projects.FindByOwner, "owner@test.io", manyprojects, []Project{pthree}, nil},
		{"FindByOwner multiple projects", Projects.FindByOwner, "owner@test.net", manyprojects, []Project{pone, ptwo}, nil},
	}
	findProjectsByIDTests = []findProjectsByIDTest{
		{"FindByID empty tables", Projects.FindByID, 42, emptytables, Project{}, sql.ErrNoRows},
		{"FindByID multiple projects none match", Projects.FindByID, 42, manyprojects, Project{}, sql.ErrNoRows},
		{"FindByID one project", Projects.FindByID, 101, oneproject, pone, nil},
		{"FindByID multiple projects", Projects.FindByID, 103, manyprojects, pthree, nil},
	}
	findProjectsByNameTests = []findProjectsByNameTest{
		{"FindByName empty tables", Projects.FindByName, "unknown", emptytables, Project{}, sql.ErrNoRows},
		{"FindByName multiple projects none match", Projects.FindByName, "unknown", manyprojects, Project{}, sql.ErrNoRows},
		{"FindByName one project", Projects.FindByName, "project one", oneproject, pone, nil},
		{"FindByName multiple projects", Projects.FindByName, "project two", manyprojects, ptwo, nil},
	}
	newProjectTests = []newProjectTest{
		{"NewProject no projects", NewProject, "project one", "first test project", "owner@test.net", 101, noprojects, pone, []Project{pone}, nil},
		{"NewProject no such owner", NewProject, "project bogus", "bogus test project", "unknown@bogus.io", 999, noprojects, Project{}, []Project{}, ErrNoSuchOwner},
		{"NewProject project exists", NewProject, "project one", "first test project", "owner@test.net", 101, oneproject, Project{}, []Project{pone}, ErrProjectExists},
		{"NewProject one project", NewProject, "project three", "third test project", "owner@test.io", 103, oneproject, pthree, []Project{pone, pthree}, nil},
		{"NewProject", NewProject, "project four", "fourth test project", "owner@test.net", 104, manyprojects, pfour, []Project{pone, ptwo, pthree, pfour}, nil},
	}
	openProjectIssueTests = []openProjectIssueTest{
		{"OpenIssue no issues", pone.OpenIssue, "issueone", "issue one", "fred@testrock.org", 1, 2001, projectissues, issueone, []Issue{issueone}, nil},
		{"OpenIssue issue exists", ptwo.OpenIssue, "issuetwo", "issue two", "barney@testrock.org", 2, 2002, projectissues, Issue{}, []Issue{issuetwo}, ErrIssueExists},
		{"OpenIssue no such reporter", pone.OpenIssue, "issueone", "issue one", "fred.c.dobbs@sierramadre.gld", 1, 2001, projectissues, Issue{}, []Issue{issuetwo}, ErrNoSuchMember},
		{"OpenIssue", ptwo.OpenIssue, "issueseven", "issue seven", "barney@testrock.org", 2, 2007, projectissues, issueseven, []Issue{issuetwo, issueseven}, nil},
		{"OpenIssue same name, different projects", ptwo.OpenIssue, "issueone", "issue one", "fred@testrock.org", 2, 2012, projectissues, issueonep2, []Issue{issuetwo, issueonep2}, nil},
	}
)

func TestContributors(t *testing.T) {
	for _, nt := range contributorstests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		cs, err := nt.fn(ctx)
		if err != nil {
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
			dropdb(db)
			continue
		}

		if !samemembers(nt.expected, cs) {
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, cs)
			dropdb(db)
			continue
		}

		dropdb(db)
	}
}

func TestFindAllProjects(t *testing.T) {
	for _, nt := range findAllProjectsTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		ps, err := nt.fn(struct{}{}, ctx)
		if err != nil {
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
			dropdb(db)
			continue
		}

		if !sameprojects(nt.expected, ps) {
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, ps)
			dropdb(db)
			continue
		}

		dropdb(db)
	}
}

func TestFindProjectsByOwner(t *testing.T) {
	for _, nt := range findProjectsByOwnerTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		ps, err := nt.fn(struct{}{}, ctx, nt.owner)
		if err != nil {
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
			dropdb(db)
			continue
		}

		if !sameprojects(nt.expected, ps) {
			t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, ps)
			dropdb(db)
			continue
		}

		dropdb(db)
	}
}

func TestFindProjectsByID(t *testing.T) {
	for _, nt := range findProjectsByIDTests {
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

func TestFindProjectsByName(t *testing.T) {
	for _, nt := range findProjectsByNameTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := nt.fn(struct{}{}, ctx, nt.name)
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

func TestNewProject(t *testing.T) {
	for _, nt := range newProjectTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)
		ids := ctx.Value("ids-chan").(chan int)

		go func() {
			ids <- nt.id
		}()

		p, err := nt.fn(ctx, nt.name, nt.desc, nt.owner)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil:
			break
		default:
			var projects Projects

			if p != nt.expected {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, p)
				break
			}

			ps, err := projects.FindAll(ctx)
			if err != nil {
				t.Errorf("%s: unexpected verification error [%+v]", nt.description, err)
				break
			}

			if !sameprojects(ps, nt.collection) {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.collection, ps)
				break
			}
		}

		dropdb(db)
	}
}

func TestOpenProjectIssue(t *testing.T) {
	for _, nt := range openProjectIssueTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)
		ids := ctx.Value("ids-chan").(chan int)

		go func() {
			ids <- nt.id
		}()

		i, err := nt.fn(ctx, nt.name, nt.desc, nt.reporter, nt.priority)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil:
			break
		default:
			var issues Issues

			if i != nt.expected {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, i)
				break
			}

			is, err := issues.FindByProject(ctx, i.project)
			if err != nil {
				t.Errorf("%s: unexpected verification error [%+v]", nt.description, err)
				break
			}

			if !sameissues(is, nt.collection) {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.collection, is)
				break
			}
		}

		dropdb(db)
	}
}

func contributors() context.Context {
	db := createdb("contributors")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create transaction to setup the database: [%+v]", err))
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
						  (102, 1006),
						  (103, 1004),
						  (103, 1005),
						  (103, 1006);`); err != nil {
		panic(fmt.Sprintf("cannot setup contributors table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func manyprojects() context.Context {
	db := createdb("manyprojects")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
						  (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES 
						  (1001, "owner@test.net"),
						  (1002, "owner@test.io");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func noprojects() context.Context {
	db := createdb("oneproject")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES (1001, "owner@test.net");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func oneproject() context.Context {
	db := createdb("oneproject")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES (101, "project one", "first test project", 1001);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO members VALUES (1001, "owner@test.net"), (1002, "owner@test.io");`); err != nil {
		panic(fmt.Sprintf("cannot setup members table: [%+v]", err))
	}

	tx.Commit()

	return ctx
}

func projectissues() context.Context {
	db := createdb("projectissues")

	ctx := context.Background()
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "ids-chan", make(chan int))

	tx, err := db.Begin()
	if err != nil {
		panic(fmt.Sprintf("cannot create a transaction to setup the database: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO projects VALUES 
	                      (101, "project one", "first test project", 1001),
						  (102, "project two", "second test project", 1001),
						  (103, "project three", "third test project", 1002);`); err != nil {
		panic(fmt.Sprintf("cannot setup projects table: [%+v]", err))
	}

	if _, err := tx.Exec(`INSERT INTO issues VALUES
	                      (2002, "issuetwo", "issue two", 2, "OPEN", 102, 1008),
						  (2011, "issueone", "issue one", 2, "OPEN", 103, 1008);`); err != nil {
		panic(fmt.Sprintf("cannot setup issues table: [%+v]", err))
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

	tx.Commit()

	return ctx
}
