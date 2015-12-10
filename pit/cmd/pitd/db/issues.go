package db

import (
	"database/sql"

	"golang.org/x/net/context"
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

// FindAll retrieves a list of all the issues in the repository.
func (Issues) FindAll(ctx context.Context) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID ORDER BY issues.ID;")
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByID retrieves the issue with the given id.
func (Issues) FindByID(ctx context.Context, id int) (Issue, error) {
	db := databaseFromContext(ctx)
	result := Issue{}

	row := db.QueryRow("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND issues.ID == $1;", id)
	if err := row.Scan(&result.id, &result.description, &result.priority, &result.status, &result.project, &result.reporter); err != nil {
		return Issue{}, err
	}

	return result, nil
}

// FindByProject retrieves a list of all the issues associated with the given project.
func (Issues) FindByProject(ctx context.Context, projectid int) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND issues.Project == $1 ORDER BY issues.ID;", projectid)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByReport retrieves a list of all the issues reported by the specified reporter.
func (Issues) FindByReporter(ctx context.Context, reporter string) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND members.Email == $1 ORDER BY issues.ID;", reporter)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByPriority retrieves a list of all the issues known to the PIT system with the given priority.
func (Issues) FindByPriority(ctx context.Context, priority int) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND issues.Priority == $1 ORDER BY issues.ID;", priority)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
}

// FindByStatus retrieves a list of all the issues know to the PIT system with the specified status.
func (Issues) FindByStatus(ctx context.Context, status string) ([]Issue, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT issues.ID, issues.Description, issues.Priority, issues.Status, issues.Project, members.Email FROM issues, members WHERE issues.Reporter == members.ID AND issues.Status == $1 ORDER BY issues.ID;", status)
	if err != nil {
		return []Issue{}, err
	}

	return collectIssues(ctx, rows)
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

		if err := rows.Scan(&issue.id, &issue.description, &issue.priority, &issue.status, &issue.project, &issue.reporter); err != nil {
			return []Issue{}, err
		}

		result = append(result, issue)
	}

	return result, nil
}
