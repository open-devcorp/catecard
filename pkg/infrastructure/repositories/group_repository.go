package repositories

import (
	"catecard/pkg/domain/entities"
	"context"
	"database/sql"
	"errors"
	"log"
)

type groupRepository struct {
	log *log.Logger
	db  *sql.DB
}

// GetByCatechistsId implements GroupRepository.
func (g *groupRepository) GetByCatechistsId(catechistId int) (int, error) {

	query := `SELECT id FROM groups WHERE catechist_id = ?`
	row := g.db.QueryRow(query, catechistId)
	var groupId int
	err := row.Scan(&groupId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No group found for the given catechist ID
		}
		g.log.Printf("Error getting group by catechist ID: %v", err)
		return 0, err
	}
	return groupId, nil

}

// Update implements GroupRepository.
func (g *groupRepository) Update(group *entities.Group) (*entities.Group, error) {
	query := `UPDATE groups SET name = ?, catechist_id = ? WHERE id = ?`
	_, err := g.db.Exec(query, group.Name, group.CatechistId, group.ID)
	if err != nil {
		g.log.Printf("Error updating group: %v", err)
		return nil, err
	}
	return group, nil
}

// DeleteById implements GroupRepository.
func (g *groupRepository) DeleteById(id int) error {
	query := `DELETE FROM groups WHERE id = ?`
	_, err := g.db.Exec(query, id)
	if err != nil {
		g.log.Printf("Error deleting group by ID: %v", err)
		return err
	}
	return nil
}

func NewGroupRepository(logger *log.Logger, db *sql.DB) GroupRepository {
	return &groupRepository{logger, db}
}

// Add implements GroupRepository.
func (g *groupRepository) Add(group *entities.Group) error {

	query := `INSERT INTO groups(name, catechist_id) VALUES (?,?)`
	result, err := g.db.Exec(query, group.Name, group.CatechistId)
	if err != nil {
		g.log.Printf("Error inserting group: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		g.log.Printf("Error getting last insert ID: %v", err)
		return err
	}
	group.ID = int(id)
	return nil
}

// Edit implements GroupRepository.
func (g *groupRepository) Edit(group *entities.Group) error {
	panic("unimplemented")
}

// GetAll implements GroupRepository.
func (g *groupRepository) GetAll() ([]*entities.Group, error) {
	query := `SELECT id, name, catechist_id FROM groups`
	rows, err := g.db.Query(query)
	if err != nil {
		g.log.Printf("Error querying groups: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []*entities.Group
	for rows.Next() {
		group := &entities.Group{}
		err := rows.Scan(&group.ID, &group.Name, &group.CatechistId)
		if err != nil {
			g.log.Printf("Error scanning group: %v", err)
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		g.log.Printf("Error iterating over group rows: %v", err)
		return nil, err
	}

	return groups, nil

}

func (g *groupRepository) GetById(id int) (*entities.Group, error) {
	query := `SELECT id, name, catechist_id FROM groups WHERE id = ?`
	row := g.db.QueryRow(query, id)

	group := &entities.Group{}
	err := row.Scan(&group.ID, &group.Name, &group.CatechistId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No group found with the given ID
		}
		g.log.Printf("Error getting group by ID: %v", err)
		return nil, err
	}

	return group, nil
}

func (r *groupRepository) Get(id int) (*GroupInfo, error) {
	const q = `
    SELECT
      g.id,
      g.name,
      g.catechist_id,
      u.id       AS user_id,
      u.username,
      u.email,
      u.role,
      COUNT(c.id) AS catechumens_count
    FROM groups g
    LEFT JOIN users u       ON u.id = g.catechist_id
    LEFT JOIN catechumens c ON c.group_id = g.id
    WHERE g.id = $1
    GROUP BY g.id, g.name, g.catechist_id, u.id, u.username, u.email, u.role;
    `
	row := r.db.QueryRowContext(context.Background(), q, id)

	grp := &entities.Group{}
	usr := &entities.User{}
	var count int

	if err := row.Scan(
		&grp.ID,
		&grp.Name,
		&grp.CatechistId,
		&usr.ID,
		&usr.Username,
		&usr.Email,
		&usr.Role,
		&count,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // no existe
		}
		r.log.Printf("GetById scan error: %v", err)
		return nil, err
	}

	// Si no hay catequista (LEFT JOIN), usr.ID puede quedar 0: maneja nil
	var cate *entities.User
	if usr.ID != 0 {
		cate = usr
	}

	return &GroupInfo{
		Group:          grp,
		Catechist:      cate,
		CatechumenSize: count,
	}, nil
}
