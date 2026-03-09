package postgres

import (
	"database/sql"

	"article/internal/models"

	"github.com/google/uuid"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	row := r.DB.QueryRow(
		"SELECT id, name, email, password FROM users WHERE email = $1",
		email,
	)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	_, err := r.DB.Exec("INSERT INTO users (id, name, password, email) VALUES ($1, $2, $3, $4);", user.ID, user.Name, user.Password, user.Email)
	return err
}

func (r *AuthRepository) CheckUserId(userId uuid.UUID) (bool, error) {
	var exists bool

	err := r.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
		userId,
	).Scan(&exists)

	return exists, err
}
