package lib

import (
	"net/http"
	"regexp"
)

func ValidateUserPayload(payload *User) *ApiResponse {

	if payload.Role != "user" && payload.Role != "admin" {
		return &ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "unkown role field",
			Response: Res{
				Error: "unknown role field",
				Data:  nil,
			},
		}
	}

	if len(payload.Phone) != 10 {
		return &ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "Phone number must be 10 character long",
			Response: Res{
				Error: "invalid phone",
				Data:  nil,
			},
		}
	}

	if len(payload.Password) <= 7 {
		return &ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "password cannot be less the 8 character long",
			Response: Res{
				Error: "invalid phone",
				Data:  nil,
			},
		}
	}

	if len(payload.Username) <= 2 {
		return &ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "username must be greater then 2 character long",
			Response: Res{
				Error: "invalid phone",
				Data:  nil,
			},
		}
	}

	expr, err := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$`)
	if err != nil {
		return &ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "regular expr error",
			Response: Res{
				Error: err,
				Data:  nil,
			},
		}
	}

	if !expr.Match([]byte(payload.Email)) {
		return &ApiResponse{
			Status:  http.StatusBadRequest,
			Message: "provied email is invalid",
			Response: Res{
				Error: "invalid email",
				Data:  nil,
			},
		}
	}

	return nil
}
