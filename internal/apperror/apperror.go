package apperror

type AppError struct {
	Message string
	Code    int
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}
