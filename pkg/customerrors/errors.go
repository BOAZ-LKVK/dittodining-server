package customerrors

var _ error = &ApplicationError{}

type ApplicationError struct {
	Code int
	Err  error
}

func (e *ApplicationError) Error() string {
	return e.Err.Error()
}
