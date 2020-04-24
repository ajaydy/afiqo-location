package helpers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type (
	Error struct {
		Err        error
		StatusCode int
		Message    string
	}
)

func (e *Error) Error() string {
	return e.Message
}

func ErrorWrap(err error, prefix, suffix, message string, status int) *Error {
	logger.Err.Errorf("error : %v : %v : %v", prefix, suffix, err)
	return &Error{
		Err:        errors.Wrapf(err, "%s/%s", prefix, suffix),
		Message:    message,
		StatusCode: status,
	}
}

func ErrorResponse(w http.ResponseWriter, message string, status int) {
	resp := Response{
		Data: nil,
		BaseResponse: BaseResponse{
			Errors: []string{
				message,
			},
		},
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		return
	}
}

const (
	InternalServerError      = "Internal Server Error"
	BadRequestMessage        = "Bad Request"
	UnauthorizedMessage      = "Unauthorized"
	ForbiddenMessage         = "Forbidden Message"
	IncorrectEmailMessage    = "Incorrect Email"
	IncorrectPasswordMessage = "Incorrect Password"
)
