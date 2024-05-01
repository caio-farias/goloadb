package components

type ServerError struct {
	status int
	mssg   string
}

func NewServerError(status int, mssg string) *ServerError {
	return &ServerError{
		status: status,
		mssg:   mssg,
	}
}

func (s *ServerError) Error() string {
	return s.mssg
}
