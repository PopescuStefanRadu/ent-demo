package request

type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUser struct {
	Id       int    `uri:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetFilteredUsers struct {
	IdsIn []int `json:"ids_in"`
}
