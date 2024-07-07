package lib

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// comparing users password
func ComparePassword(hash, password, salt string) bool {
	passHash := HashGenerator(password, salt)
	return passHash == hash
}

func ValidateUserPayload(payload *UserPayload) *ApiResponse {

	if payload.AccountType != "seller" && payload.AccountType != "buyer" {
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

func ValidateLoginPayload(payload UserLoginPayload) *ApiResponse {
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

// CancelHttpRequestOperation waits for either a 5-second delay or a cancellation signal from the context.
// If the delay completes first, it returns http.StatusOK.
// If the context is canceled first, it returns http.StatusRequestTimeout.
func CancleHttpRequestOpearation(ctx context.Context) int {
	select {
	case <-time.After(5 * time.Second):
		return http.StatusRequestTimeout

	case <-ctx.Done():
		return http.StatusOK
	}
}

// Custom Response Writer to write data to http Response
func WriteResponse[T any](w http.ResponseWriter, status int, payload T) {
	w.WriteHeader(status)
	byteBuffer, _ := PasreString(&payload)
	w.Write(byteBuffer)
}

// parses string into json
func ParseJson[T any](r *http.Request) (*T, error) {

	var data T
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// parses struct or json into string

func PasreString[T any](data *T) ([]byte, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	return buf, nil
}

func HashGenerator(password string, salt string) string {
	hasher := sha256.New()
	saltedPassword := append([]byte(password), []byte(salt)...)
	hasher.Write(saltedPassword)
	hashedword := hasher.Sum(nil)
	hexVal := hex.EncodeToString(hashedword)
	return hexVal
}

func JwtWriter(id string, email string, key string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"id":    id,
	})

	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ParseJwt(token, secret string) (*JwtData, error) {

	if secret == "" {
		return nil, errors.New("jwt signing secret not found")
	}

	sercetByte, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid sigining method")
		}

		return []byte(secret), nil

	})
	if err != nil {
		return nil, err
	}

	if values, ok := sercetByte.Claims.(jwt.MapClaims); ok && sercetByte.Valid {
		email, _ := values["email"].(string)
		id, _ := values["id"].(string)
		return &JwtData{
			Email: email,
			Id:    id,
		}, nil
	}

	return nil, errors.New("invalid token")
}
