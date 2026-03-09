package repositories

import (
	"article/internal/models"

	"github.com/google/uuid"
)

type AuthRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(u *models.User) error
	CheckUserId(id uuid.UUID) (bool, error)
}
