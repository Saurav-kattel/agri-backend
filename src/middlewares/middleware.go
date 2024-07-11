package middlewares

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sauravkattel/agri/src/lib"
	"sauravkattel/agri/src/users"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type writerWrapper struct {
	http.ResponseWriter
	Status int
}

type Middleware func(http.Handler) http.Handler

func (w *writerWrapper) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.Status = statusCode

}

func LoggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		status := &writerWrapper{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(status, r)
		log.Println(status.Status, r.Method, r.URL.Path, time.Since(start))
	})
}

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

type UsersContextKeyType string

const UsersContextKey UsersContextKeyType = "users_key"

func AuthMiddleware(db *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = godotenv.Load()
			token := r.Header.Get("auth_token")

			if token == "" || token == "undefined" {
				lib.WriteResponse(
					w,
					http.StatusBadRequest,
					lib.ApiResponse{
						Status:  http.StatusBadRequest,
						Message: "auth token not found",
						Response: lib.Res{
							Error: "Authentication not found",
						},
					},
				)
				return

			}
			key := os.Getenv("KEY")
			userData, err := lib.ParseJwt(token, key)
			if err != nil {
				lib.WriteResponse(
					w,
					http.StatusInternalServerError,
					lib.ApiResponse{
						Status:  http.StatusInternalServerError,
						Message: "error occured while parsing jwt token",
						Response: lib.Res{
							Error: err.Error(),
						},
					},
				)
				return
			}

			dbUser, err := users.GetUsersByUserId(db, userData.Id)
			if err != nil && err == sql.ErrNoRows {
				lib.WriteResponse(
					w,
					http.StatusNotFound,
					lib.ApiResponse{
						Status:  http.StatusNotFound,
						Message: "user with this id not found",
						Response: lib.Res{
							Error: err.Error(),
						},
					},
				)
				return
			}
			if err != nil {
				lib.WriteResponse(
					w,
					http.StatusInternalServerError,
					lib.ApiResponse{
						Status:  http.StatusInternalServerError,
						Message: "error occured while fetching users",
						Response: lib.Res{
							Error: err.Error(),
						},
					},
				)
				return

			}

			ctx := context.WithValue(r.Context(), UsersContextKey, dbUser)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
