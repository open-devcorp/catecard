package repositories

import (
	"catecard/pkg/domain/entities"
	"context"
	"database/sql"
	"log"
	"time"
)

type userRepository struct {
	log *log.Logger
	db  *sql.DB
}

// GetById implements UserRepository.
func (r *userRepository) GetById(id int) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT id, username, full_name, password, role FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Username, &user.FullName, &user.Password, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Printf("User not found with ID: %d", id)
			return nil, nil
		}
		r.log.Println("Error getting user by ID:", err)
		return nil, err
	}

	return user, nil
}

// GetAll implements usecases.UserRepository.
func (r *userRepository) GetAll() ([]*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, username, full_name, password, role FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}

		err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Password, &user.Role)
		if err != nil {
			r.log.Println("Error scanning user:", err)
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		r.log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetUser(username, password string) *entities.User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Validate by username AND password, and retrieve full_name as well
	query := `SELECT id, username, full_name, password, role FROM users WHERE username = ? AND password = ?`
	row := r.db.QueryRowContext(ctx, query, username, password)

	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Username, &user.FullName, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Printf("User not found: %s", username)
			return nil
		}
		r.log.Println("Error getting user:", err)
		return nil
	}

	return user
}

// SaveUser implements usecases.UserRepository.
func (r *userRepository) SaveUser(user *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Println("Error starting transaction:", err)
		return nil, err
	}
	defer tx.Rollback()

	// Insert username and full_name only (email removed from schema)
	query := `INSERT INTO users (username, full_name, password, role) VALUES (?, ?, ?, ?) RETURNING id`

	err = tx.QueryRowContext(ctx, query, user.Username, user.FullName, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		r.log.SetFlags(log.LstdFlags | log.Llongfile)
		r.log.Println("Error saving user:", err)
		r.log.SetFlags(log.LstdFlags)
		return nil, err
	}

	// Convert timestamps to strings

	if err = tx.Commit(); err != nil {
		r.log.Println("Error committing transaction:", err)
		return nil, err
	}

	return user, nil
}

func (r *userRepository) DeleteUserById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Println("Error deleting user:", err)
		return err
	}

	return nil
}
func NewUserRepository(log *log.Logger, database *sql.DB) UserRepository {
	return &userRepository{log, database}
}
