package users

import (
	"sauravkattel/agri/src/lib"

	"github.com/jmoiron/sqlx"
)

func CreateUser(db *sqlx.DB, payload *lib.UserPayload, hash string) (*string, error) {
	var id string
	err := db.QueryRowx(
		"INSERT INTO users(email,password,phone,username,account_type,address) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		payload.Email,
		hash,
		payload.Phone,
		payload.Username,
		payload.AccountType,
		payload.Address,
	).Scan(&id)
	if err != nil {

		return nil, err
	}
	return &id, nil
}

func GetUsersByUserId(db *sqlx.DB, id string) (*lib.User, error) {
	data := lib.User{}
	err := db.QueryRowx("SELECT * FROM users WHERE id = $1", id).StructScan(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func GetUsersByEmail(db *sqlx.DB, email string) (*lib.User, error) {
	data := lib.User{}
	err := db.QueryRowx("SELECT * FROM users WHERE email = $1", email).StructScan(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func GetUsersByUserName(db *sqlx.DB, username string) (*lib.User, error) {
	data := lib.User{}
	err := db.QueryRowx("SELECT * FROM users WHERE username = $1", username).StructScan(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
