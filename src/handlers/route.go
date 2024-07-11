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

	authFunc := middlewares.AuthMiddleware(db)

	authMiddleware := middlewares.CreateStack(
		authFunc,
	)

	mux.Handle("/api/v1/user/register", unAuthMiddleWare(RegisterUserHandler(db, ctx)))
	mux.Handle("/api/v1/user/login", unAuthMiddleWare(LoginUserHandler(db, ctx)))

	mux.Handle("/api/v1/product/add", authMiddleware(AddProduct(db)))
	return mux
}
