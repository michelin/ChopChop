package internal

// ErrNilParameter is an error meaning a parameter is nil.
type ErrNilParameter struct {
	Name string
}

func (e ErrNilParameter) Error() string {
	return "parameter " + e.Name + " is nil"
}
