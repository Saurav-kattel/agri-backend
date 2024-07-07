package handlers

import (
	"context"
	"net/http"
	"sauravkattel/agri/src/middlewares"

	"github.com/jmoiron/sqlx"
)

func GetRoutes(db *sqlx.DB, ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()

	unAuthMiddleWare := middlewares.CreateStack(
		middlewares.LoggerMiddleWare,
	)

	mux.Handle("/api/v1/user/register", unAuthMiddleWare(RegisterUserHandler(db, ctx)))
	mux.Handle("/api/v1/user/login", unAuthMiddleWare(LoginUserHandler(db, ctx)))

	return mux
}
