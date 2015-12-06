package db

import (
	"database/sql"

	"golang.org/x/net/context"
)

// Project is a project managed by the PIT system and owned by a specific member of the project team.
type Project struct {
	id          int
	name        string
	description string
	owner       string
}

// Projects is the collection of all the projects managed by the PIT system.
type Projects struct{}

// FindAll retrieves a list of all the projects in the repository.
func (Projects) FindAll(ctx context.Context) ([]Project, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT projects.ID, projects.Name, projects.Description, members.Email FROM projects, members WHERE projects.Owner == members.ID ORDER BY projects.ID")
	if err != nil {
		return []Project{}, err
	}

	return collectProjects(ctx, rows)
}

// FindByOwner retrieves a list of all projects owned by the specified owner
func (Projects) FindByOwner(ctx context.Context, owner string) ([]Project, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT projects.ID, projects.Name, projects.Description, members.Email FROM projects, members WHERE projects.Owner == members.ID AND members.Email == $1 ORDER BY projects.ID;", owner)
	if err != nil {
		return []Project{}, err
	}

	return collectProjects(ctx, rows)
}

// FindByID retrieves the project with the given ID
func (Projects) FindByID(ctx context.Context, id int) (Project, error) {
	db := databaseFromContext(ctx)
	result := Project{}

	row := db.QueryRow("SELECT projects.ID, projects.Name, projects.Description, members.Email FROM projects, members WHERE projects.Owner == members.ID AND projects.ID == $1;", id)
	if err := row.Scan(&result.id, &result.name, &result.description, &result.owner); err != nil {
		return Project{}, err
	}

	return result, nil
}

// FindByName retrieves the project with the given name
func (Projects) FindByName(ctx context.Context, name string) (Project, error) {
	db := databaseFromContext(ctx)
	result := Project{}

	row := db.QueryRow("SELECT projects.ID, projects.Name, projects.Description, members.Email FROM projects, members WHERE projects.Owner == members.ID AND projects.Name == $1;", name)
	if err := row.Scan(&result.id, &result.name, &result.description, &result.owner); err != nil {
		return Project{}, err
	}

	return result, nil
}

// Contributors retrieves a list of project team members contributing to the project.
func (p Project) Contributors(ctx context.Context) ([]Member, error) {
	db := databaseFromContext(ctx)

	rows, err := db.Query("SELECT members.ID, members.Email FROM members FULL JOIN contributors ON (members.ID == contributors.MID) WHERE contributors.PID == $1 ORDER BY members.ID;", p.id)
	if err != nil {
		return []Member{}, err
	}

	return collectMembers(ctx, rows)
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
