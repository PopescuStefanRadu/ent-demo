package request

type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserURI struct {
	ID int `binding:"required" uri:"id"`
}

type UpdateUserBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetFilteredUsers struct {
	IdsIn []int `json:"ids_in"`
}
