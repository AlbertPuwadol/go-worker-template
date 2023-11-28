package apperror

type AppErrorCode string

const (
	KirinError AppErrorCode = "KIRIN_ERROR"
)

type AppError struct {
	Message     string
	Description string
	Code        AppErrorCode
}

func (err AppError) Error() string {
	return err.Message
}

func NewError(message, description string, code AppErrorCode) AppError {
	return AppError{
		Message:     message,
		Description: description,
		Code:        code,
	}
}
