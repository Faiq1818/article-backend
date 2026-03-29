package auths

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(req requesttype.LoginRequest) (string, error) {
	// query get user data from req.email
	user, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user: %v", err)
			return "", &pkg.AppError{
				Message: "User not found",
				Code:    400,
				Err:     err,
			}
		}
		log.Printf("Database scan error: %v", err)
		return "", &pkg.AppError{
			Message: "Database error",
			Code:    500,
			Err:     err,
		}
	}

	// check and compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", &pkg.AppError{
			Message: "Incorrect password",
			Code:    400,
			Err:     err,
		}
	}

	// making the jwt
	key := []byte(os.Getenv("JWT_SECRET"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": user.ID,
			"role":    "admin",
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		})

	token, err := t.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return "", &pkg.AppError{
			Message: "Failed to generate token",
			Code:    500,
			Err:     err,
		}
	}

	s.Logger.Info("Login success")

	// response
	return token, nil
}
