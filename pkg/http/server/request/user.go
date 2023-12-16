package request

type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserURI struct {
	Id int `json:"id" uri:"id" binding:"required"`
}

type UpdateUserBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetFilteredUsers struct {
	IdsIn []int `json:"ids_in"`
}
