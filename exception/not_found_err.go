package exception

type NotFoundError struct {
	// disini akan mengikuti kontrak interface error
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	// nah disini kita menggunakan return value nya adalah struct NotFoundError
	return NotFoundError{Error: error}
}