package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sauravkattel/agri/src/lib"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		payload    lib.UserPayload
		wantStatus int
		wantBody   string
	}{
		{
			name:   "valid response",
			status: http.StatusOK,
			payload: lib.UserPayload{
				Username:    "johndoe",
				Email:       "john.doe@example.com",
				Password:    "password123",
				Phone:       "1234567890",
				AccountType: "user",
				Address:     "123 Main St",
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"username":"johndoe","email":"john.doe@example.com","password":"password123","phone":"1234567890","account_type":"user","address":"123 Main St"}`,
		},
		{
			name:       "empty payload",
			status:     http.StatusBadRequest,
			payload:    lib.UserPayload{},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"username":"","email":"","password":"","phone":"","account_type":"","address":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			lib.WriteResponse(rr, tt.status, tt.payload)
			assert.Equal(t, tt.wantStatus, rr.Code, "status code mismatch")
			assert.JSONEq(t, tt.wantBody, rr.Body.String(), "response body mismatch")
		})
	}
}

func TestCancleHttpRequestOpearation(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "passed request",
			args: args{
				ctx: func() context.Context {
					ctx, cancle := context.WithCancel(context.Background())
					defer cancle()
					return ctx
				}(),
			},
			want: 200,
		}, {
			name: "cancle request",
			args: args{
				ctx: func() context.Context {
					ctx := context.Background()
					return ctx
				}(),
			},
			want: http.StatusRequestTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.CancleHttpRequestOpearation(tt.args.ctx); got != tt.want {
				t.Errorf("CancleHttpRequestOpearation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {

	type args struct {
		password string
		salt     string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "run hash test",
			args: args{
				password: "myvar",
				salt:     "hi",
			},
			want: "47de6b8943d4c679ee1c5c927383b1b324f54c7f863c8fe4a532d97ab7328ac2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.HashGenerator(tt.args.password, tt.args.salt)
			if got != tt.want {
				t.Errorf("HashGenereator() want = %+v, got = %+v", tt.want, got)
			}
		})
	}

}

func TestJwtWriter(t *testing.T) {

	type args struct {
		email string
		id    string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "run hash test",
			args: args{
				id:    "myvar",
				email: "hi",
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImhpIiwiaWQiOiJteXZhciJ9.gAoIG9FJIQM28idWKRAZqc8t6wGNyYlNUt5q17yFuNk",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lib.JwtWriter(tt.args.id, tt.args.email, "saurav")
			if (err != nil) != tt.wantErr {
				t.Errorf("error want %+v got %+v", tt.wantErr, err)
			}
			if got != tt.want {
				t.Errorf("HashGenereator() want = %+v, got = %+v", tt.want, got)
			}
		})
	}

}
