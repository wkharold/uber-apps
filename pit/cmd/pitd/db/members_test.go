package db

import (
	"database/sql"
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
	ctxfn       func() context.Context
	expected    []Member
	err         error
}

type findMemberByEmailTest struct {
	description string
	email       string
	ctxfn       func() context.Context
	expected    Member
	err         error
}

type findMemberByIDTest struct {
	description string
	id          int
	ctxfn       func() context.Context
	expected    Member
	err         error
}

type joinProjectTest struct {
	description string
	member      Member
	project     Project
	ctxfn       func() context.Context
	projects    []Project
	err         error
}

type memberIssuesTest struct {
	description string
	fn          func(context.Context) ([]Issue, error)
	ctxfn       func() context.Context
	expected    []Issue
	err         error
}

type memberWatchTest struct {
	description string
	member      Member
	issue       Issue
	ctxfn       func() context.Context
	watching    []Issue
	err         error
}

type newMemberTest struct {
	description string
	email       string
	id          int
	ctxfn       func() context.Context
	expected    Member
	collection  []Member
	err         error
}

var (
	bob   = Member{id: 1003, email: "bob@members.com"}
	carol = Member{id: 1004, email: "carol@members.com"}
	ted   = Member{id: 1005, email: "ted@members.com"}
	alice = Member{id: 1006, email: "alice@members.com"}
	wilma = Member{id: 1009, email: "wilma@testrock.org"}
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
		{"FindAll no members", emptytables, []Member{}, nil},
		{"FindAll one member", onemember, []Member{bob}, nil},
		{"FindAll members", manymembers, []Member{carol, ted, alice}, nil},
	}
	findMemberByEmailTests = []findMemberByEmailTest{
		{"FindByEmail empty tables", "bob@members.com", emptytables, Member{}, sql.ErrNoRows},
		{"FindByEmail many members no match", "fred.c.dobbs@members.com", manymembers, Member{}, sql.ErrNoRows},
		{"FindByEmail one member", "bob@members.com", onemember, bob, nil},
		{"FindByEmail members", "ted@members.com", manymembers, ted, nil},
	}
	findMemberByIDTests = []findMemberByIDTest{
		{"FindByID empty tables", 1003, emptytables, Member{}, sql.ErrNoRows},
		{"FindByID many members no match", 2001, manymembers, Member{}, sql.ErrNoRows},
		{"FindByID one member", 1003, onemember, bob, nil},
		{"FindByID members", 1005, manymembers, ted, nil},
	}
	joinProjectTests = []joinProjectTest{
		{"Join non existent project", bob, Project{}, emptytables, []Project{}, ErrNoSuchProject},
		{"Join first project", bob, pone, oneproject, []Project{pone}, nil},
		{"Join another project", alice, pone, contributors, []Project{pone, ptwo, pthree}, nil},
		{"Join already contributing", alice, ptwo, contributors, []Project{ptwo, pthree}, nil},
	}
	memberWatchTests = []memberWatchTest{
		{"Watch non existent issue", bob, Issue{}, alltheissues, []Issue{}, ErrNoSuchIssue},
		{"Watch first issue", wilma, issueone, alltheissues, []Issue{issueone}, nil},
		{"Watch another issue", alice, issueone, alltheissues, []Issue{issueone, issuefive}, nil},
	}
	newMemberTests = []newMemberTest{
		{"NewMember empty tables", "bob@members.com", 1003, emptytables, bob, []Member{bob}, nil},
		{"NewMember member exists", "bob@members.com", 1003, onemember, Member{}, []Member{}, ErrMemberExists},
		{"NewMember", "alice@members.com", 1006, onemember, alice, []Member{bob, alice}, nil},
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

func TestMemberJoin(t *testing.T) {
	for _, nt := range joinProjectTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		err := nt.member.Join(ctx, nt.project)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil:
			break
		default:
			projects, err := nt.member.ContributesTo(ctx)
			if err != nil {
				t.Errorf("%s: cannot retrieve project member contributes to: [%+v]", nt.description, err)
				break
			}

			if !sameprojects(projects, nt.projects) {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.projects, projects)
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

func TestMemberWatch(t *testing.T) {
	for _, nt := range memberWatchTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		err := nt.member.Watch(ctx, nt.issue)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil:
			break
		default:
			watching, err := nt.member.Watching(ctx)
			if err != nil {
				t.Errorf("%s: cannot retrieve issues being watched: [%+v]", nt.description, err)
				break
			}

			if !sameissues(watching, nt.watching) {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.watching, watching)
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

		ms, err := FindAllMembers(ctx)
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

func TestFindMemberByEmail(t *testing.T) {
	for _, nt := range findMemberByEmailTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := FindMemberByEmail(ctx, nt.email)
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

func TestFindMemberByID(t *testing.T) {
	for _, nt := range findMemberByIDTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)

		p, err := FindMemberByID(ctx, nt.id)
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

func TestNewMembers(t *testing.T) {
	for _, nt := range newMemberTests {
		ctx := nt.ctxfn()
		db := ctx.Value("database").(*sql.DB)
		ids := ctx.Value("ids-chan").(chan int)

		go func() {
			ids <- nt.id
		}()

		m, err := NewMember(ctx, nt.email)
		switch {
		case err != nil && err != nt.err:
			t.Errorf("%s: unexpected error [%+v]", nt.description, err)
		case err != nil:
			break
		default:
			if m != nt.expected {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.expected, m)
				break
			}

			ms, err := FindAllMembers(ctx)
			if err != nil {
				t.Errorf("%s: unexpected verification error [%+v]", nt.description, err)
				break
			}

			if !samemembers(ms, nt.collection) {
				t.Errorf("%s: expected %+v, got %+v", nt.description, nt.collection, ms)
				break
			}
		}

		dropdb(db)
	}
}
