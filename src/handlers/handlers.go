package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sauravkattel/agri/src/lib"
	"sauravkattel/agri/src/middlewares"
	"sauravkattel/agri/src/product"
	"sauravkattel/agri/src/users"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func RegisterUserHandler(db *sqlx.DB, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		SALT := os.Getenv("SALT")
		KEY := os.Getenv("KEY")

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
		hashPassword := lib.HashGenerator(userPayload.Password, SALT)
		id, err := users.CreateUser(db, userPayload, hashPassword)

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
		token, err := lib.JwtWriter(*id, userPayload.Email, KEY)

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

func LoginUserHandler(db *sqlx.DB, ctx context.Context) http.HandlerFunc {
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
		userPayload, err := lib.ParseJson[lib.UserLoginPayload](r)
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

		errors := lib.ValidateLoginPayload(*userPayload)
		if errors != nil {
			lib.WriteResponse(w, http.StatusBadRequest, errors)
			return
		}

		data, err := users.GetUsersByEmail(db, userPayload.Email)
		if err != nil {
			lib.WriteResponse(
				w,
				http.StatusInternalServerError,
				lib.ApiResponse{
					Status:  http.StatusInternalServerError,
					Message: "error occured while fetching users by email",
					Response: lib.Res{
						Error: err.Error(),
					},
				},
			)
			return
		}

		SALT := os.Getenv("SALT")
		KEY := os.Getenv("KEY")

		log.Println(SALT, KEY)
		ok := lib.ComparePassword(data.Password, userPayload.Password, SALT)
		log.Println(ok, data.Password, userPayload.Password)
		if !ok {
			lib.WriteResponse(
				w,
				http.StatusBadRequest,
				lib.ApiResponse{
					Status:  http.StatusBadRequest,
					Message: "invlaid password",
					Response: lib.Res{
						Error: "password didnot match",
					},
				},
			)

			return
		}

		token, err := lib.JwtWriter(data.Id, data.Email, KEY)

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

func AddProduct(db *sqlx.DB) http.HandlerFunc {
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
		productPayload, err := lib.ParseJson[lib.ProductPayload](r)
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

		userData, ok := r.Context().Value(middlewares.UsersContextKey).(*lib.User)
		if !ok {
			lib.WriteResponse(w, http.StatusInternalServerError, lib.ApiResponse{
				Status:  http.StatusInternalServerError,
				Message: "error occured while fetching users",
				Response: lib.Res{
					Error: "cannot parse data from req context",
				},
			})
			return

		}

		charSeq := lib.GetRandomCharSequence()
		productSlug := strings.Join(strings.Split(productPayload.Product.Name, " "), "-") + "-" + charSeq

		err = product.AddProduct(db, productPayload.Product, productPayload.Attrib, userData.Id, productSlug)
		if err != nil {
			lib.WriteResponse(w, http.StatusInternalServerError, lib.ApiResponse{
				Status:  http.StatusInternalServerError,
				Message: "error occured while inserting products",
				Response: lib.Res{
					Error: err.Error(),
				},
			})
			return
		}

		lib.WriteResponse(
			w,
			http.StatusOK,
			lib.ApiResponse{
				Status:  http.StatusOK,
				Message: "success",
				Response: lib.Res{
					Data: "succrss",
				},
			},
		)

	}
}
