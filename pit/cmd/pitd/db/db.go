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
	db  *sql.DB
	IDs chan<- int
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
	db, err := sql.Open("ql", "memory://pit.db")
	if err != nil {
		panic(fmt.Sprintf("cannot create database instance: [%+v]", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("database ping failed: [%+v]", err))
	}

	if err = mkTables(); err != nil {
		panic(fmt.Sprintf("table creation failed: [%+v]", err))
	}

	go nextID(IDs)
}

// FindAll retrieves a list of all the issues in the repository.
func (Issues) FindAll(ctx context.Context) ([]Issue, error) {
	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Report = members.ID;")
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByProject retrieves a list of all the issues associated with the given project.
func (Issues) FindByProject(ctx context.Context, projectid int) ([]Issue, error) {
	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID AND issues.Project = ?;", projectid)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByReport retrieves a list of all the issues reported by the specified reporter.
func (Issues) FindByReporter(ctx context.Context, reporter string) ([]Issue, error) {
	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID AND members.Email = ?;", reporter)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByPriority retrieves a list of all the issues known to the PIT system with the given priority.
func (Issues) FindByPriority(ctx context.Context, priority int) ([]Issue, error) {
	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID AND issues.Priority = ?;", priority)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByStatus retrieves a list of all the issues know to the PIT system with the specified status.
func (Issues) FindByStatus(ctx context.Context, status string) ([]Issue, error) {
	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter = members.ID AND issues.Status = ?;", status)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindAll retrieves a list of all the projects in the repository.
func (Projects) FindAll(ctx context.Context) ([]Project, error) {
	rows, err := db.Query("SELECT projects.ID, projects.Name, projects.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID")
	if err != nil {
		return []Project{}, err
	}

	return collectProjects(ctx, rows)
}

// FindByOwner retrieves a list of all projects owned by the specified owner
func (Projects) FindByOwner(ctx context.Context, owner string) ([]Project, error) {
	rows, err := db.Query("SELECT projects.ID, projects.Name, projects.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID AND members.Email = ?;", owner)
	if err != nil {
		return []Project{}, err
	}

	return collectProjects(ctx, rows)
}

// FindByID retrieves the project with the given ID
func (Projects) FindByID(ctx context.Context, id int) (Project, error) {
	result := Project{}

	row := db.QueryRow("SELECT projects.ID, projects.Name, projejcts.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID AND projects.ID = ?;", id)
	if err := row.Scan(&result.id, &result.name, &result.description, &result.owner); err != nil {
		return Project{}, err
	}

	return result, nil
}

// FindByName retrieves the project with the given name
func (Projects) FindByName(ctx context.Context, name string) (Project, error) {
	result := Project{}

	row := db.QueryRow("SELECT projects.ID, projects.Name, projejcts.Description, members.Email FROM projects, members WHERE projects.Owner = members.ID AND projects.Name = ?;", name)
	if err := row.Scan(&result.id, &result.name, &result.description, &result.owner); err != nil {
		return Project{}, err
	}

	return result, nil
}

func collectIssues(ctx context.Context, rows *sql.Rows) ([]Issue, error) {
	result := []Issue{}

	for rows.Next() {
		issue := Issue{}

		if err = rows.Scan(&issue.id, &issue.description, &issue.priority, &issue.status, &issue.project, &issue.reporter); err != nil {
			return []Issue{}, err
		}

		result = append(result, issue)
	}

	return result, nil
}

func collectProjects(ctx context.Context, rows *sql.Rows) ([]Project, error) {
	result := []Project{}

	for rows.Next() {
		project := Project{}

		if err = rows.Scan(&project.id, &project.name, &project.description, &project.owner); err != nil {
			return []Project{}, err
		}

		result = append(result, project)
	}

	return result, nil
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

	if _, err = tx.Exec("CREATE TABLE issues (ID int, Description string, Priority int, Status string, Project int,  Reporter int);"); err != nil {
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
