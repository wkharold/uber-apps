package db

import (
	"database/sql"
	"testing"

	"golang.org/x/net/context"
)

type addMemberTest struct {
	description string
	project     Project
	member      Member
	ctxfn       func() context.Context
	postmembers []Member
	err         error
}

type contributorstest struct {
	description string
	fn          func(context.Context) ([]Member, error)
	ctxfn       func() context.Context
	expected    []Member
	err         error
}

type findAllProjectsTest struct {
	description string
	ctxfn       func() context.Context
	expected    []Project
	err         error
}

type findProjectsByOwnerTest struct {
	description string
	owner       string
	ctxfn       func() context.Context
	expected    []Project
	err         error
}

type findProjectByIDTest struct {
	description string
	id          int
	ctxfn       func() context.Context
	expected    Project
	err         error
}

type findProjectByNameTest struct {
	description string
	name        string
	ctxfn       func() context.Context
	expected    Project
	err         error
}

type newProjectTest struct {
	description string
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

type removeMemberTest struct {
	description string
	project     Project
	member      Member
	ctxfn       func() context.Context
	postmembers []Member
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
	addMemberTests = []addMemberTest{
		{"AddMember no such member", pone, Member{id: 42, email: "fred.c.dobbs@sierramadre.gld"}, emptytables, []Member{}, ErrNoSuchMember},
		{"AddMember first contributor", pone, bob, contributors, []Member{bob}, nil},
		{"AddMember duplicate contributor", pthree, ted, contributors, []Member{carol, ted, alice}, nil},
		{"AddMember additional contributor", ptwo, bob, contributors, []Member{bob, alice}, nil},
	}
	contributorstests = []contributorstest{
		{"Contributors no contributors", pone.Contributors, contributors, []Member{}, nil},
		{"Contributors single contributor", ptwo.Contributors, contributors, []Member{alice}, nil},
		{"Contributors multiple", pthree.Contributors, contributors, []Member{carol, ted, alice}, nil},
	}
	findAllProjectsTests = []findAllProjectsTest{
		{"FindAll from empty tables", emptytables, []Project{}, nil},
		{"FindAll one project", oneproject, []Project{pone}, nil},
		{"FindAll multiple projects", manyprojects, []Project{pone, ptwo, pthree}, nil},
	}
	findProjectsByOwnerTests = []findProjectsByOwnerTest{
		{"FindByOwner from empty tables", "owner@test.net", emptytables, []Project{}, nil},
		{"FindByOwner one project no match", "owner@test.org", oneproject, []Project{}, nil},
		{"FindByOwner multiple projects no match", "owner@test.com", manyprojects, []Project{}, nil},
		{"FindByOwner one project", "owner@test.net", oneproject, []Project{pone}, nil},
		{"FindByOwner multiple projects one match", "owner@test.io", manyprojects, []Project{pthree}, nil},
		{"FindByOwner multiple projects", "owner@test.net", manyprojects, []Project{pone, ptwo}, nil},
	}
	findProjectByIDTests = []findProjectByIDTest{
		{"FindByID empty tables", 42, emptytables, Project{}, sql.ErrNoRows},
		{"FindByID multiple projects none match", 42, manyprojects, Project{}, sql.ErrNoRows},
		{"FindByID one project", 101, oneproject, pone, nil},
		{"FindByID multiple projects", 103, manyprojects, pthree, nil},
	}
	findProjectByNameTests = []findProjectByNameTest{
		{"FindByName empty tables", "unknown", emptytables, Project{}, sql.ErrNoRows},
		{"FindByName multiple projects none match", "unknown", manyprojects, Project{}, sql.ErrNoRows},
		{"FindByName one project", "project one", oneproject, pone, nil},
		{"FindByName multiple projects", "project two", manyprojects, ptwo, nil},
	}
	newProjectTests = []newProjectTest{
		{"NewProject no projects", "project one", "first test project", "owner@test.net", 101, noprojects, pone, []Project{pone}, nil},
		{"NewProject no such owner", "project bogus", "bogus test project", "unknown@bogus.io", 999, noprojects, Project{}, []Project{}, ErrNoSuchOwner},
		{"NewProject project exists", "project one", "first test project", "owner@test.net", 101, oneproject, Project{}, []Project{pone}, ErrProjectExists},
		{"NewProject one project", "project three", "third test project", "owner@test.io", 103, oneproject, pthree, []Project{pone, pthree}, nil},
		{"NewProject", "project four", "fourth test project", "owner@test.net", 104, manyprojects, pfour, []Project{pone, ptwo, pthree, pfour}, nil},
	}
	openProjectIssueTests = []openProjectIssueTest{
		{"OpenIssue no issues", pone.OpenIssue, "issueone", "issue one", "fred@testrock.org", 1, 2001, projectissues, issueone, []Issue{issueone}, nil},
		{"OpenIssue issue exists", ptwo.OpenIssue, "issuetwo", "issue two", "barney@testrock.org", 2, 2002, projectissues, Issue{}, []Issue{issuetwo}, ErrIssueExists},
		{"OpenIssue no such reporter", pone.OpenIssue, "issueone", "issue one", "fred.c.dobbs@sierramadre.gld", 1, 2001, projectissues, Issue{}, []Issue{issuetwo}, ErrNoSuchMember},
		{"OpenIssue", ptwo.OpenIssue, "issueseven", "issue seven", "barney@testrock.org", 2, 2007, projectissues, issueseven, []Issue{issuetwo, issueseven}, nil},
		{"OpenIssue same name, different projects", ptwo.OpenIssue, "issueone", "issue one", "fred@testrock.org", 2, 2012, projectissues, issueonep2, []Issue{issuetwo, issueonep2}, nil},
	}
	removeMemberTests = []removeMemberTest{
		{"RemoveMember no such member", pone, Member{id: 42, email: "fred.c.dobbs@sierramadre.gld"}, emptytables, []Member{}, ErrNoSuchMember},
		{"RemoveMember only member", ptwo, alice, contributors, []Member{}, nil},
		{"RemoveMember remove a member", pthree, ted, contributors, []Member{carol, alice}, nil},
		{"RemoveMember member with assignment", pthree, carol, alltheissues, []Member{carol, ted, alice}, ErrMemberHasIssues},
	}
)

func TestAddMember(t *testing.T) {
	for _, nt := range addMemberTests {
		ctx := nt.ctxfn()
		db, ok := ctx.Value("database").(*sql.DB)
		if !ok {
			t.Fatalf("%s: no database in context")
		}

		err := nt.project.AddMember(ctx, nt.member)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil:
			break
		default:
			members, err := nt.project.Contributors(ctx)
			if err != nil {
				t.Errorf("%s: cannot retrieve contributors: [%+v]", nt.description, err)
				break
			}

			if !samemembers(members, nt.postmembers) {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.postmembers, members)
			}
		}

		dropdb(db)
	}
}

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

		ps, err := FindAllProjects(ctx)
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

		ps, err := FindProjectsByOwner(ctx, nt.owner)
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

func TestFindProjectByID(t *testing.T) {
	for _, nt := range findProjectByIDTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := FindProjectByID(ctx, nt.id)
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

func TestFindProjectByName(t *testing.T) {
	for _, nt := range findProjectByNameTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := FindProjectByName(ctx, nt.name)
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

		p, err := NewProject(ctx, nt.name, nt.desc, nt.owner)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil:
			break
		default:
			if p != nt.expected {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, p)
				break
			}

			ps, err := FindAllProjects(ctx)
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
			if i != nt.expected {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, i)
				break
			}

			is, err := FindIssuesByProject(ctx, i.project)
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

func TestRemoveMember(t *testing.T) {
	for _, nt := range removeMemberTests {
		ctx := nt.ctxfn()
		db, ok := ctx.Value("database").(*sql.DB)
		if !ok {
			t.Fatalf("%s: no database in context")
		}

		err := nt.project.RemoveMember(ctx, nt.member.id)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil:
			break
		default:
			members, err := nt.project.Contributors(ctx)
			if err != nil {
				t.Errorf("%s: cannot retrieve contributors: [%+v]", nt.description, err)
				break
			}

			if !samemembers(members, nt.postmembers) {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.postmembers, members)
			}
		}

		dropdb(db)
	}
}
