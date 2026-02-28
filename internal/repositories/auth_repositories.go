package repositories

import "article/internal/models"

type AuthRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(u *models.User) error
}
