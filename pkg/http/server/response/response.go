package response

type Response[T any] struct {
	Result T                  `json:"result,omitempty"`
	Errors map[string][]Error `json:"errors,omitempty"`
}

type Error struct {
	Cause   error  `json:"-"`
	Path    string `json:"-"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (err *Error) Error() string {
	if err == nil {
		return ""
	}

	return err.Message
}

func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}

	return err.Cause
}
