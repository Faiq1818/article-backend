package auths

import (
	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(req requesttype.LoginRequest) error {
	// query get user data from req.email
	userData := h.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ($1);", req.Email)

	var id, name, email, password string
	err := userData.Scan(&id, &name, &email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user: %v", err)
			return &pkg.AppError{
				Message: "User tidak ditemukan",
				Code:    400,
				Err:     err,
			}
		}
		log.Printf("Database scan error: %v", err)
		return &pkg.AppError{
			Message: "Database error",
			Code:    500,
			Err:     err,
		}
	}

	// check and compare password
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password))
	if err != nil {
		return &pkg.AppError{
			Message: "Password salah",
			Code:    400,
			Err:     err,
		}
	}

	// making the jwt
	key := []byte(os.Getenv("JWT_SECRET"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"role": "admin",
			"name": "Faiq",
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		})

	s, err := t.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return &pkg.AppError{
			Message: "Error saat membuat token",
			Code:    500,
			Err:     err,
		}
	}
	_ = s

	// response
	return nil
}
