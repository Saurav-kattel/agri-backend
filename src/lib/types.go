package lib

type Res struct {
	Error any
	Data  any
}

type ApiResponse struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Response Res
}

type User struct {
	Id        string `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
	Password  string `json:"passowrd" db:"password"`
	Role      string `json:"role" db:"role"`
}
