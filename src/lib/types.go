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

type JwtData struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type UserLoginPayload struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"passowrd" db:"password"`
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

type Product struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"dec"`
}

type Attrib struct {
	Price      float32 `json:"price" db:"price"`
	Quantity   float32 `json:"quantity" db:"quantity"`
	Status     string  `json:"status" db:"status"`
	Product_id *string `json:"product_id" db:"product_id"`
	Slug       string  `json:"slug" db:"slug"`
}

type ProductDetails struct {
	Id          string  `json:"id" db:"products.id"`
	Name        string  `json:"name" db:"products.name"`
	Description string  `json:"description" db:"products.description"`
	CreatedAt   *string `json:"created_at" db:"products.created_at"`
	Price       float32 `json:"price" db:"pa.price"`
	Quantity    float32 `json:"quantity" db:"pa.quantity"`
	Status      string  `json:"status" db:"pa.status"`
	Attrib_Id   string  `json:"attrib_id" db:"pa.id as attrib_id"`
	Product_id  *string `json:"product_id" db:"pa.products_id"`
	Slug        string  `json:"slug" db:"pa.slug"`
}
