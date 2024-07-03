package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"sauravkattel/agri/src/lib"
	"sauravkattel/agri/src/users"
	"time"

	"github.com/jmoiron/sqlx"
)

func RegisterUserHandler(db *sqlx.DB, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			lib.WriteResponse(w, http.StatusMethodNotAllowed, lib.ApiResponse{
				Status:  http.StatusMethodNotAllowed,
				Message: "cannot use other http method then POST",
				Response: lib.Res{
					Error: "invalid http request method",
				},
			})
			return
		}

		// parsing req.body into json or struct
		userPayload, err := lib.ParseJson[lib.UserPayload](r)
		if err != nil {
			lib.WriteResponse(w, http.StatusInternalServerError, lib.ApiResponse{
				Status:  http.StatusInternalServerError,
				Message: "error occured while parsing req.body",
				Response: lib.Res{
					Error: err.Error(),
				},
			})
			return

		}

		// validating users josn data
		errResponse := lib.ValidateUserPayload(userPayload)
		if errResponse != nil {
			lib.WriteResponse(w, http.StatusBadRequest, errResponse)
			return
		}

		// finding dupicate users with the email
		_, err = users.GetUsersByEmail(db, userPayload.Email)

		// if err occures but error is not the no rows find then send error response
		if err != nil && err != sql.ErrNoRows {
			lib.WriteResponse(w, http.StatusInternalServerError, lib.ApiResponse{
				Status:  http.StatusInternalServerError,
				Message: "error occured getting user data by email",
				Response: lib.Res{
					Error: err.Error(),
				},
			})
			return
		}

		// if err doesn't occur then that indicates that the user alredy exists so return error
		if err == nil {
			lib.WriteResponse(w, http.StatusUnauthorized, lib.ApiResponse{
				Status:  http.StatusUnauthorized,
				Message: "users with this email already exists",
				Response: lib.Res{
					Error: "duplicate email error",
				},
			})
			return

		}

		// create a new user
		id, err := users.CreateUser(db, userPayload, "")
		if err != nil {
			lib.WriteResponse(w, http.StatusInternalServerError, lib.ApiResponse{
				Status:  http.StatusInternalServerError,
				Message: "error occured while creating new user",
				Response: lib.Res{
					Error: err.Error(),
				},
			})
			return
		}

		//signing jwt token
		salt := os.Getenv("SALT")
		token, err := lib.JwtWriter(*id, userPayload.Email, salt)

		if err != nil {
			lib.WriteResponse(w, http.StatusInternalServerError, lib.ApiResponse{
				Status:  http.StatusInternalServerError,
				Message: "error occured while creating jwt token",
				Response: lib.Res{
					Error: err.Error(),
				},
			})
			return
		}

		cookie := &http.Cookie{
			Name:  "auth_token_cookie",
			Value: token,
			Path:  "/",

			Expires:  time.Now().Add(time.Hour * 24 * 10),
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)

		lib.WriteResponse(w, http.StatusOK, lib.ApiResponse{
			Status:  http.StatusOK,
			Message: "registered successfully",
		})

	}
}
