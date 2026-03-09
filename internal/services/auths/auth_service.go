package auths

import (
	pkg "article/internal/pkg"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Service) CheckUserAuthorization(userInfo jwt.MapClaims) (string, error) {
	userIDStr, ok := userInfo["user_id"].(string)
	if !ok {
		return "", &pkg.AppError{
			Message: "Invalid user_id",
			Code:    400,
		}
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", &pkg.AppError{
			Message: "Invalid UUID",
			Code:    400,
			Err:     err,
		}
	}

	exist, err := s.Repo.CheckUserId(userID)
	if err != nil {
		s.Logger.Error("Database error", "err", err)
		return "", &pkg.AppError{
			Message: "Database error",
			Code:    500,
			Err:     err,
		}
	}

	if !exist {
		return "", &pkg.AppError{
			Message: "User is not exist",
			Code:    403,
			Err:     err,
		}
	}

	return "User authorized", nil
}
