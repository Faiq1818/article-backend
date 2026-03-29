package pkg

import (
	"github.com/lib/pq"
	"net/http"
)

func ParsePostgresError(err error) (int, string) {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505": // unique_violation
			return http.StatusConflict, "Data already exists (duplicate)."
		case "23503": // foreign_key_violation
			return http.StatusBadRequest, "Referenced data not found."
		case "23502": // not_null_violation
			return http.StatusBadRequest, "Required field is empty."
		case "22001": // string_data_right_truncation
			return http.StatusBadRequest, "Input text is too long."
		}
	}

	return http.StatusInternalServerError, "Internal server error occurred."
}
