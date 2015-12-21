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

// NewProject creates a new project with the specified name, description, and owner
func NewProject(ctx context.Context, name, description, owner string) (Project, error) {
	db := databaseFromContext(ctx)
	ids := idsChanFromContext(ctx)
	id := <-ids

	created := false

	tx, err := db.Begin()
	if err != nil {
		return Project{}, err
	}
	defer func() {
		switch created {
		case true:
			tx.Commit()
		case false:
			tx.Rollback()
		}
	}()

	var memail, pname string
	var mid int

	err = tx.QueryRow("SELECT ID, Email FROM members WHERE Email == $1", owner).Scan(&mid, &memail)
	if err != nil {
		return Project{}, ErrNoSuchOwner
	}

	err = tx.QueryRow("SELECT Name FROM projects WHERE Name == $1", name).Scan(&pname)
	if err != sql.ErrNoRows {
		return Project{}, ErrProjectExists
	}

	_, err = tx.Exec("INSERT INTO projects VALUES ($1, $2, $3, $4)", id, name, description, mid)
	if err != nil {
		return Project{}, err
	}

	created = true

	return Project{id: id, name: name, description: description, owner: owner}, nil
}

// AddMember adds the specified member to the project's list of contributors.
func (p Project) AddMember(ctx context.Context, member Member) error {
	db := databaseFromContext(ctx)

	added := false

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		switch added {
		case true:
			tx.Commit()
		case false:
			tx.Rollback()
		}
	}()

	var mid int
	err = tx.QueryRow("SELECT ID FROM members WHERE ID == $1", member.id).Scan(&mid)
	if err != nil {
		return ErrNoSuchMember
	}

	rows, err := tx.Query("SELECT members.ID, members.Email FROM members FULL JOIN contributors ON (members.ID == contributors.MID) WHERE contributors.PID == $1 ORDER BY members.ID;", p.id)
	if err != nil {
		return err
	}

	members, err := collectMembers(ctx, rows)
	if err != nil {
		return err
	}

	for _, m := range members {
		if m == member {
			// Already a member, nothing to do
			return nil
		}
	}

	_, err = tx.Exec("INSERT INTO contributors VALUES ($1, $2)", p.id, member.id)
	if err != nil {
		return err
	}

	added = true
	return nil
}

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

// OpenIssue creates a new issue associated with this project.
func (p Project) OpenIssue(ctx context.Context, name, description, reporter string, priority int) (Issue, error) {
	db := databaseFromContext(ctx)
	ids := idsChanFromContext(ctx)
	id := <-ids

	created := false

	tx, err := db.Begin()
	if err != nil {
		return Issue{}, err
	}
	defer func() {
		switch created {
		case true:
			tx.Commit()
		case false:
			tx.Rollback()
		}
	}()

	var iname, memail string
	var mid int

	err = tx.QueryRow("SELECT ID, Email FROM members WHERE Email == $1", reporter).Scan(&mid, &memail)
	if err != nil {
		return Issue{}, ErrNoSuchMember
	}

	err = tx.QueryRow("SELECT Name FROM issues WHERE Name == $1 AND Project == $2", name, p.id).Scan(&iname)
	if err != sql.ErrNoRows {
		return Issue{}, ErrIssueExists
	}

	_, err = tx.Exec("INSERT INTO issues VALUES($1, $2, $3, $4, $5, $6, $7)", id, name, description, priority, Open, p.id, mid)
	if err != nil {
		return Issue{}, err
	}

	created = true

	return Issue{id: id, name: name, description: description, priority: priority, status: Open, project: p.id, reporter: reporter}, nil
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
