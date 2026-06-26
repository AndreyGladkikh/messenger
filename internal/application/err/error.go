package err

type ErrorWithDetails struct {
	Message string
	Details map[string]any
}

func (e *ErrorWithDetails) Error() string {
	return e.Message
}
