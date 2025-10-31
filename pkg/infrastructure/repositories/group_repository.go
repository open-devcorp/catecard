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

// UpdateLimitGroup implements GroupRepository.
func (g *groupRepository) UpdateLimitGroup(groupId int) error {
	// Obtiene el grupo
	grp, err := g.GetById(groupId)
	if err != nil {
		g.log.Printf("Error fetching group for UpdateLimitGroup: %v", err)
		return err
	}
	if grp == nil {
		g.log.Printf("UpdateLimitGroup: group %d not found", groupId)
		return sql.ErrNoRows
	}

	// Aumenta el límite en 1 (si se requiere otro incremento, cambiar aquí)
	grp.LimitCatechumens = grp.LimitCatechumens + 1

	// Persiste el cambio
	_, err = g.Update(grp)
	if err != nil {
		g.log.Printf("Error updating group limit for id %d: %v", groupId, err)
		return err
	}

	return nil
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
	query := `UPDATE groups SET name = ?, catechist_id = ?, limit_catechumens= ? WHERE id = ?`
	_, err := g.db.Exec(query, group.Name, group.CatechistId, group.LimitCatechumens, group.ID)
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

func (g *groupRepository) Add(group *entities.Group) error {
	query := `INSERT INTO groups(name, catechist_id, limit_catechumens) VALUES (?,?,?)`
	result, err := g.db.Exec(query, group.Name, group.CatechistId, group.LimitCatechumens)
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

// GetAll implements GroupRepository.
func (g *groupRepository) GetAll() ([]*entities.Group, error) {
	query := `SELECT id, name, catechist_id,limit_catechumens FROM groups`
	rows, err := g.db.Query(query)
	if err != nil {
		g.log.Printf("Error querying groups: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []*entities.Group
	for rows.Next() {
		group := &entities.Group{}
		err := rows.Scan(&group.ID, &group.Name, &group.CatechistId, &group.LimitCatechumens)
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
	query := `SELECT id, name, catechist_id, limit_catechumens FROM groups WHERE id = ?`
	row := g.db.QueryRow(query, id)

	group := &entities.Group{}
	err := row.Scan(&group.ID, &group.Name, &group.CatechistId, &group.LimitCatechumens)
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
				g.limit_catechumens,
				u.id       AS user_id,
				u.username,
				u.full_name,
				u.role,
				COUNT(c.id) AS catechumens_count
			FROM groups g
			LEFT JOIN users u       ON u.id = g.catechist_id
			LEFT JOIN catechumens c ON c.group_id = g.id
			WHERE g.id = $1
			GROUP BY g.id, g.name, g.catechist_id, g.limit_catechumens, u.id, u.username, u.full_name, u.role;
			`
	row := r.db.QueryRowContext(context.Background(), q, id)

	grp := &entities.Group{}
	usr := &entities.User{}
	var count int

	if err := row.Scan(
		&grp.ID,
		&grp.Name,
		&grp.CatechistId,
		&grp.LimitCatechumens,
		&usr.ID,
		&usr.Username,
		&usr.FullName,
		&usr.Role,
		&count,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // no existe
		}
		r.log.Printf("GetById scan error: %v", err)
		return nil, err
	}

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
