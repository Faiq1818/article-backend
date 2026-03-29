package auths

import (
	models "article/internal/models"
	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(req requesttype.RegisterRequest) error {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("Error when hasing password from new user", "err", err)
		return &pkg.AppError{
			Message: "failed to process password",
			Code:    500,
			Err:     err,
		}
	}

	// insert to db
	user := &models.User{
		Email:    req.Email,
		Name:     req.Name,
		ID:       uuid.New(),
		Password: string(hashedPassword),
	}

	err = s.Repo.CreateUser(user)
	if err != nil {
		s.Logger.Error("Error when push new user to db", "err", err)
		return &pkg.AppError{
			Message: "Failed to create account, ensure email is unique",
			Code:    500,
			Err:     err,
		}
	}

	// response
	return nil
}
