package db

import (
	"database/sql"

	"golang.org/x/net/context"
)

// Issue is an issue reported by a member of a project team.
type Issue struct {
	id          int
	name        string
	description string
	priority    int
	status      string
	project     int
	reporter    string
}

// Description retrieves the description of an issue.
func (i Issue) Description() string {
	return i.description
}

// ID retrieves an issue's id.
func (i Issue) ID() int {
	return i.id
}

// Name retrieves the name of an issue.
func (i Issue) Name() string {
	return i.name
}

// Priority retrieves an issue's priority.
func (i Issue) Priority() int {
	return i.priority
}

// Project returns the id of the project the issue belongs to.
func (i Issue) Project() int {
	return i.project
}

// Reporter retrieves the email address of the team member who reported the issue.
func (i Issue) Reporter() string {
	return i.reporter
}

// Status retrieves the status of the issue.
func (i Issue) Status() string {
	return i.status
}

// FindAll retrieves a list of all the issues in the repository.
func FindAllIssues(ctx context.Context) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Name, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID ORDER BY issues.ID;")
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByID retrieves the issue with the given id.
func FindIssueByID(ctx context.Context, id int) (Issue, error) {
	db := databaseFromContext(ctx)
	result := Issue{}

	row := db.QueryRow("SELECT issues.ID, issues.Name, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND issues.ID == $1;", id)
	if err := row.Scan(&result.id, &result.name, &result.description, &result.priority, &result.status, &result.project, &result.reporter); err != nil {
		return Issue{}, err
	}

	return result, nil
}

// FindByProject retrieves a list of all the issues associated with the given project.
func FindIssuesByProject(ctx context.Context, projectid int) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Name, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND issues.Project == $1 ORDER BY issues.ID;", projectid)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByReport retrieves a list of all the issues reported by the specified reporter.
func FindIssuesByReporter(ctx context.Context, reporter string) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Name, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND members.Email == $1 ORDER BY issues.ID;", reporter)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByPriority retrieves a list of all the issues known to the PIT system with the given priority.
func FindIssuesByPriority(ctx context.Context, priority int) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Name, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND issues.Priority == $1 ORDER BY issues.ID;", priority)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByStatus retrieves a list of all the issues know to the PIT system with the specified status.
func FindIssuesByStatus(ctx context.Context, status string) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Name, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND issues.Status == $1 ORDER BY issues.ID;", status)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// Assign adds the given Member to the list of Members assigned to this issue. The Member must exist and be a
// contributor to the Project to which this issue belongs.
func (i Issue) Assign(ctx context.Context, m Member) error {
	db := databaseFromContext(ctx)

	assigned := false

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		switch assigned {
		case true:
			tx.Commit()
		case false:
			tx.Rollback()
		}
	}()

	var mid int

	err = tx.QueryRow("SELECT ID FROM members WHERE ID == $1", m.id).Scan(&mid)
	if err != nil {
		return ErrNoSuchMember
	}

	// get the list of issues this member is currently assigned
	rows, err := tx.Query(`
	SELECT I.IID, I.Name, I.Description, I.Priority, I.Status, I.Project, I.Reporter
	FROM (SELECT issues.ID AS IID, issues.Name AS Name, issues.Description AS Description, issues.Priority AS Priority, issues.Status AS Status, issues.Project AS Project, members.Email AS Reporter
	      FROM issues, members
		  WHERE issues.Reporter == members.ID
		  ORDER BY IID) AS I
	FULL JOIN assignments ON (I.IID == assignments.IID)
    WHERE assignments.MID == $1
	ORDER BY I.IID
	`, m.id)
	if err != nil {
		return err
	}

	assignments, err := collectIssues(ctx, rows)

	for _, assignment := range assignments {
		if assignment == i {
			// already assigned, nothing to do
			return nil
		}
	}

	// get a list of projects this member contributes to
	rows, err = tx.Query("SELECT members.ID, members.Email FROM members FULL JOIN contributors ON (members.ID == contributors.MID) WHERE contributors.PID == $1 ORDER BY members.ID;", i.project)
	if err != nil {
		return err
	}

	contributors, err := collectMembers(ctx, rows)
	if err != nil {
		return err
	}

	for _, contributor := range contributors {
		if contributor == m {
			goto AssignMember // Yes, it's a goto. The code is cleaner than it would be without it.
		}
	}

	return ErrNonContributingMember

AssignMember:
	_, err = tx.Exec("INSERT INTO assignments VALUES ($1, $2)", m.id, i.id)
	if err != nil {
		return err
	}

	assigned = true
	return nil
}

// Assigned retrieves a list of project team members who are assigned to the issue.
func (i Issue) Assigned(ctx context.Context) ([]Member, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT members.ID, members.Email FROM members FULL JOIN assignments ON (members.ID == assignments.MID) WHERE assignments.IID == $1 ORDER BY members.ID;", i.id)
	if err != nil {
		return []Member{}, err
	}

	return collectMembers(ctx, rows)
}

// Watching retrieves a list of project team members who are watching the issue.
func (i Issue) Watching(ctx context.Context) ([]Member, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT members.ID, members.Email FROM members FULL JOIN watchers ON (members.ID == watchers.MID) WHERE watchers.IID == $1 ORDER BY members.ID;", i.id)
	if err != nil {
		return []Member{}, err
	}

	return collectMembers(ctx, rows)
}

func collectIssues(ctx context.Context, rows *sql.Rows) ([]Issue, error) {
	result := []Issue{}

	for rows.Next() {
		issue := Issue{}

		if err := rows.Scan(&issue.id, &issue.name, &issue.description, &issue.priority, &issue.status, &issue.project, &issue.reporter); err != nil {
			return []Issue{}, err
		}

		result = append(result, issue)
	}

	rows.Close()

	return result, nil
}
