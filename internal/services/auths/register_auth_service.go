package auths

import (
	"log"

	pkg "article/internal/pkg"
	requesttype "article/internal/request_type"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Register(req requesttype.RegisterRequest) error {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return &pkg.AppError{
			Message: "failed to process password",
			Code:    500,
			Err:     err,
		}
	}

	// insert to db
	u := uuid.New()
	_, err = h.DB.Exec("INSERT INTO users (id, name, password, email) VALUES ($1, $2, $3, $4);", u, req.Name, hashedPassword, req.Email)
	if err != nil {
		log.Println(err)
		return &pkg.AppError{
			Message: "Akun gagal dibuat, pastikan email unik",
			Code:    500,
			Err:     err,
		}
	}

	// response
	return nil
}
