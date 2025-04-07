package deepl

import "strconv"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (self Error) Error() string {
	return "code: " + strconv.Itoa(self.Code) + ", message: " + self.Message
}
