package users

import (
	"sauravkattel/agri/src/lib"

	"github.com/jmoiron/sqlx"
)

func CreateUser(db *sqlx.DB, payload *lib.UserPayload, hash string, role string) error {
	_, err := db.Exec(
		"INSERT INTO users(first_name,last_name,email,password,phone,username,role) VALUES ($1,$2,$3,$4,$5,$6,$7)",
		payload.FirstName,
		payload.LastName,
		payload.Email,
		hash,
		payload.Phone,
		payload.Username,
		role,
	)

	return err
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
