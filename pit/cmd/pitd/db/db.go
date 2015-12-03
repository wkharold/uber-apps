// Package db provides the pit database implementation
package db

import (
	"database/sql"
	"fmt"

	"golang.org/x/net/context"

	_ "github.com/cznic/ql"
)

const (
	Closed   = "CLOSED"
	Open     = "OPEN"
	Returned = "RETURNED"
)

var (
	qldb *sql.DB
	IDs  chan<- int
)

// Issue is an issue reported by a member of a project team.
type Issue struct {
	id          int
	description string
	priority    int
	status      string
	project     int
	reporter    string
}

// Issues is the collection of all of the reported issues known to the PIT system.
type Issues struct{}

// Member is a member of a project team.
type Member struct {
	id    int
	email string
}

// Members is the collection of all the project team members known to the PIT system.
type Members struct{}

// Project is a project managed by the PIT system and owned by a specific member of the project team.
type Project struct {
	id          int
	name        string
	description string
	owner       string
}

// Projects is the collection of all the projects managed by the PIT system.
type Projects struct{}

func init() {
	qldb, err := sql.Open("ql", "memory://pit.db")
	if err != nil {
		panic(fmt.Sprintf("cannot create database instance: [%+v]", err))
	}

	if err = qldb.Ping(); err != nil {
		panic(fmt.Sprintf("database ping failed: [%+v]", err))
	}

	if err = mkTables(qldb); err != nil {
		panic(fmt.Sprintf("table creation failed: [%+v]", err))
	}

	go nextID(IDs)
}

// FindAll retrieves a list of all the issues in the repository.
func (Issues) FindAll(ctx context.Context) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID;")
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByID retrieves the issue with the given id.
func (Issues) FindByID(ctx context.Context, id int) (Issue, error) {
	db := databaseFromContext(ctx)
	result := Issue{}

	row := db.QueryRow("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID AND issues.ID = ?;", id)
	if err := row.Scan(&result.id, &result.description, &result.priority, &result.status, &result.project, &result.reporter); err != nil {
		return Issue{}, err
	}

	return result, nil
}

// FindByProject retrieves a list of all the issues associated with the given project.
func (Issues) FindByProject(ctx context.Context, projectid int) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID AND issues.Project = ?;", projectid)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByReport retrieves a list of all the issues reported by the specified reporter.
func (Issues) FindByReporter(ctx context.Context, reporter string) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID AND members.Email = ?;", reporter)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByPriority retrieves a list of all the issues known to the PIT system with the given priority.
func (Issues) FindByPriority(ctx context.Context, priority int) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID AND issues.Priority = ?;", priority)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByStatus retrieves a list of all the issues know to the PIT system with the specified status.
func (Issues) FindByStatus(ctx context.Context, status string) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID AND issues.Status = ?;", status)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// Assigned retrieves a list of project team members who are assigned to the issue.
func (i Issue) Assigned(ctx context.Context) ([]Member, error) {
	return []Member{}, nil
}

// Watching retrieves a list of project team members who are watching the issue.
func (i Issue) Watching(ctx context.Context) ([]Member, error) {
	return []Member{}, nil
}

// FindAll retreives a list of all the project team members in the repository.
func (Members) FindAll(ctx context.Context) ([]Member, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT members.ID, members.Email FROM members")
	if err != nil {
		return []Member{}, err
	}

	return collectMembers(ctx, rows)
}

// FindByEmail retrieves the project team member with the given email address.
func (Members) FindByEmail(ctx context.Context, email string) (Member, error) {
	db := databaseFromContext(ctx)
	result := Member{}

	err := db.QueryRow("SELECT members.ID, members.Email FROM members WHERE members.Email = ?", email).Scan(&result.id, &result.email)
	if err != nil {
		return Member{}, err
	}

	return result, nil
}

// FindByID retrieves the project team member with the given id.
func (Members) FindByID(ctx context.Context, memberid int) (Member, error) {
	db := databaseFromContext(ctx)
	result := Member{}

	err := db.QueryRow("SELECT members.ID, members.Email FROM members WHERE members.ID= ?", memberid).Scan(&result.id, &result.email)
	if err != nil {
		return Member{}, err
	}

	return result, nil
}

// Assignments retrieves a list of all the issues to which the proejct team member has been assigned.
func (m Member) Assignments(ctx context.Context) ([]Issue, error) {
	return []Issue{}, nil
}

// ContributesTo retrieves a list of all the projects to which the project team member contributes.
func (m Member) ContributesTo(ctx context.Context) ([]Project, error) {
	return []Project{}, nil
}

// Watching retrieves a list of all the issue the project team member is watching.
func (m Member) Watching(ctx context.Context) ([]Issue, error) {
	return []Issue{}, nil
}

// FindAll retrieves a list of all the projects in the repository.
func (Projects) FindAll(ctx context.Context) ([]Project, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT projects.ID, projects.Name, projects.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID")
	if err != nil {
		return []Project{}, err
	}

	return collectProjects(ctx, rows)
}

// FindByOwner retrieves a list of all projects owned by the specified owner
func (Projects) FindByOwner(ctx context.Context, owner string) ([]Project, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT projects.ID, projects.Name, projects.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID AND members.Email = ?;", owner)
	if err != nil {
		return []Project{}, err
	}

	return collectProjects(ctx, rows)
}

// FindByID retrieves the project with the given ID
func (Projects) FindByID(ctx context.Context, id int) (Project, error) {
	db := databaseFromContext(ctx)
	result := Project{}

	row := db.QueryRow("SELECT projects.ID, projects.Name, projejcts.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID AND projects.ID = ?;", id)
	if err := row.Scan(&result.id, &result.name, &result.description, &result.owner); err != nil {
		return Project{}, err
	}

	return result, nil
}

// FindByName retrieves the project with the given name
func (Projects) FindByName(ctx context.Context, name string) (Project, error) {
	db := databaseFromContext(ctx)
	result := Project{}

	row := db.QueryRow("SELECT projects.ID, projects.Name, projejcts.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID AND projects.Name = ?;", name)
	if err := row.Scan(&result.id, &result.name, &result.description, &result.owner); err != nil {
		return Project{}, err
	}

	return result, nil
}

// Contributors retrieves a list of project team members contributing to the project.
func (p Project) Contributors(ctx context.Context) ([]Member, error) {
	return []Member{}, nil
}

func collectIssues(ctx context.Context, rows *sql.Rows) ([]Issue, error) {
	result := []Issue{}

	for rows.Next() {
		issue := Issue{}

		if err := rows.Scan(&issue.id, &issue.description, &issue.priority, &issue.status, &issue.project, &issue.reporter); err != nil {
			return []Issue{}, err
		}

		result = append(result, issue)
	}

	return result, nil
}

func collectMembers(ctx context.Context, rows *sql.Rows) ([]Member, error) {
	result := []Member{}

	for rows.Next() {
		member := Member{}

		if err := rows.Scan(&member.id, &member.email); err != nil {
			return []Member{}, err
		}

		result = append(result, member)
	}

	return result, nil
}

func collectProjects(ctx context.Context, rows *sql.Rows) ([]Project, error) {
	result := []Project{}

	for rows.Next() {
		project := Project{}

		if err := rows.Scan(&project.id, &project.name, &project.description, &project.owner); err != nil {
			return []Project{}, err
		}

		result = append(result, project)
	}

	return result, nil
}

func databaseFromContext(ctx context.Context) *sql.DB {
	result := ctx.Value("database").(*sql.DB)
	if result != nil {
		return result
	}
	return qldb
}

func mkTables(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Entity tables: projects, issues, members
	if _, err = tx.Exec("CREATE TABLE projects (ID int, Name string, Description string, Owner int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE issues (ID int, Description string, Priority int, Status string, Project int, Reporter int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE members (ID int, Email string);"); err != nil {
		tx.Rollback()
		return err
	}

	// Association tables: contributors, assignments, watchers
	if _, err = tx.Exec("CREATE TABLE contributors (PID int, MID int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE assignments (MID int, IID int);"); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec("CREATE TABLE watchers (MID int, IID int);"); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func nextID(ch chan<- int) {
	next := 100

	for {
		ch <- next
		next++
	}
}
