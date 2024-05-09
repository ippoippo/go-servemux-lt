package demo

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	// Would normally be other useful things here
}

func ErrWithMsg(msg string) *ErrorResponse {
	return &ErrorResponse{
		ErrorMessage: msg,
	}
}
