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
	Id          string  `json:"id" db:"id"`
	Username    string  `json:"username" db:"username"`
	Email       string  `json:"email" db:"email"`
	Phone       string  `json:"phone" db:"phone"`
	Password    string  `json:"passowrd" db:"password"`
	AccountType string  `json:"account_type" db:"account_type"`
	Address     string  `json:"address" db:"address"`
	CreatedAt   *string `json:"created_at" db:"created_at"`
}

type UserPayload struct {
	Username    string `json:"username" db:"username"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	Phone       string `json:"phone" db:"phone"`
	AccountType string `json:"account_type" db:"account_type"`
	Address     string `json:"address" db:"address"`
}
