package handlers

import (
	"net/http"
	"sauravkattel/agri/src/lib"
	"sauravkattel/agri/src/middlewares"
	"sauravkattel/agri/src/product"
	"strings"

	"github.com/jmoiron/sqlx"
)

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

func DeleteProduct(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			lib.WriteResponse(w, http.StatusMethodNotAllowed, lib.ApiResponse{
				Status:  http.StatusMethodNotAllowed,
				Message: "cannot use other http method then DELETE",
				Response: lib.Res{
					Error: "invalid http request method",
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
		slug := r.URL.Query().Get("slug")
		if slug == "" || slug == "undefined" {
			lib.WriteResponse(w, http.StatusBadRequest, lib.ApiResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid slug ",
				Response: lib.Res{
					Error: "cannot parse slug from req slug is " + slug,
				},
			})
			return

		}

		err := product.DeleteProduct(db, userData.Id, slug)
		if err != nil {
			lib.WriteResponse(w, http.StatusInternalServerError, lib.ApiResponse{
				Status:  http.StatusInternalServerError,
				Message: "error occured while deleting product",
				Response: lib.Res{
					Error: err.Error(),
				},
			})
			return
		}

		lib.WriteResponse(w, http.StatusOK, lib.ApiResponse{
			Status:  http.StatusOK,
			Message: "deleted successfully",
		})
	}
}

func UpdateAttrib(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			lib.WriteResponse(w, http.StatusMethodNotAllowed, lib.ApiResponse{
				Status:  http.StatusMethodNotAllowed,
				Message: "cannot use other http method then PUT",
				Response: lib.Res{
					Error: "invalid http request method",
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

	}
}
