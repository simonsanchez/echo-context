package api

type CustomError struct {
	Code     int
	Public   string
	Internal error
}

func (c CustomError) Error() string {
	return c.Public
}

func newCustomError(code int, message string, err error) error {
	return &CustomError{
		Code:     code,
		Public:   message,
		Internal: err,
	}
}
