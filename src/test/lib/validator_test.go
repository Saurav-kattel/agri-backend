package test

import (
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

func TestHashCompare(t *testing.T) {

	type args struct {
		hash     string
		password string
		salt     string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid password test",
			args: args{
				hash:     "47de6b8943d4c679ee1c5c927383b1b324f54c7f863c8fe4a532d97ab7328ac2",
				password: "myvar",
				salt:     "hi",
			},
			want: true,
		}, {
			name: "invalid password test",
			args: args{
				hash:     "de6b8943d4c679ee1c5c927383b1b324f54c7f863c8fe4a532d97ab7328ac2",
				password: "myvar",
				salt:     "hi",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lib.ComparePassword(tt.args.hash, tt.args.password, tt.args.salt)
			if got != tt.want {
				t.Errorf("HashGenereator() want = %+v, got = %+v", tt.want, got)
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

func TestParseJwt(t *testing.T) {

	type args struct {
		token string
		key   string
	}

	tests := []struct {
		name    string
		args    args
		want    *lib.JwtData
		wantErr bool
	}{
		{
			name: "run jwt parse test",
			args: args{
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImhpIiwiaWQiOiJteXZhciJ9.gAoIG9FJIQM28idWKRAZqc8t6wGNyYlNUt5q17yFuNk",
				key:   "saurav",
			},
			want: &lib.JwtData{
				Email: "hi",
				Id:    "myvar",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lib.ParseJwt(tt.args.token, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("want err %+v got err %+v", tt.wantErr, err)
			}

			assert.NotNil(t, got, "returned value was nil")
			assert.Equalf(t, got.Email, tt.want.Email, "email didnot match want %s got %s", tt.want.Email, got.Email)
			assert.Equalf(t, got.Id, tt.want.Id, "id didnot match want %s got %s", tt.want.Id, got.Id)

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
