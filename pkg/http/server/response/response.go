package response

type Response[T any] struct {
	Result T                `json:"result,omitempty"`
	Errors map[string]Error `json:"errors,omitempty"`
}

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
