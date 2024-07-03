package handlers

import (
	"context"
	"net/http"
	"sauravkattel/agri/src/middlewares"

	"github.com/jmoiron/sqlx"
)

func GetRoutes(db *sqlx.DB, ctx context.Context) *http.ServeMux {
	mux := http.NewServeMux()

	authMiddleWare := middlewares.CreateStack(
		middlewares.LoggerMiddleWare,
	)

	mux.Handle("/api/v1/user/register", authMiddleWare(RegisterUserHandler(db, ctx)))

	return mux
}
