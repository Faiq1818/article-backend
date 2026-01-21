package utils

import (
    "net/http"
    "github.com/lib/pq" // atau pgx
)

// ParsePostgresError mengubah error DB menjadi status code dan pesan user
func ParsePostgresError(err error) (int, string) {
    // Cek apakah ini error dari driver Postgres (lib/pq)
    if pqErr, ok := err.(*pq.Error); ok {
        switch pqErr.Code {
        case "23505": // unique_violation
            return http.StatusConflict, "Data sudah ada (duplikat)."
        case "23503": // foreign_key_violation
            return http.StatusBadRequest, "Data referensi tidak ditemukan."
        case "23502": // not_null_violation
            return http.StatusBadRequest, "Ada data wajib yang kosong."
        case "22001": // string_data_right_truncation
            return http.StatusBadRequest, "Input teks terlalu panjang."
        }
    }
    
    // Default error (koneksi putus, sintaks salah, dll)
    return http.StatusInternalServerError, "Terjadi kesalahan internal pada server."
}
