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

// GetAll implements usecases.UserRepository.
func (r *userRepository) GetAll() ([]*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, username,email, password, role FROM users ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Println("Error querying users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		var createdAt, updatedAt time.Time

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &createdAt, &updatedAt)
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

	query := `SELECT id, username, password, role FROM users WHERE username = ?`
	row := r.db.QueryRowContext(ctx, query, username)

	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
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

	// SQLite query with correct syntax
	query := `INSERT INTO users (username, password,email, role) VALUES (?, ?,?, ?) RETURNING id`

	err = tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)
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

func NewUserRepository(log *log.Logger, database *sql.DB) UserRepository {
	return &userRepository{log, database}
}
