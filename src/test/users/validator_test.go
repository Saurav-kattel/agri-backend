package test

import (
	"net/http"
	"reflect"
	"sauravkattel/agri/src/lib"
	"testing"
)

func TestValidateUserPayload(t *testing.T) {
	type args struct {
		payload *lib.User
	}
	tests := []struct {
		name string
		args args
		want *lib.ApiResponse
	}{
		{
			name: "invalid email",
			args: args{
				payload: &lib.User{
					Id:        "123",
					Username:  "Sam",
					Password:  "1234567890",
					FirstName: "s",
					LastName:  "23",
					Email:     "kattelsaurav.com",
					Role:      "user",
					Phone:     "1234567890",
				},
			},
			want: &lib.ApiResponse{
				Status:  http.StatusBadRequest,
				Message: "provied email is invalid",
				Response: lib.Res{
					Error: "invalid email",
					Data:  nil,
				}},
		},
		{
			name: "invalid password",
			args: args{
				payload: &lib.User{
					Id:        "123",
					Username:  "Sam",
					Password:  "1230",
					Phone:     "1234567890",
					FirstName: "s",
					LastName:  "23",
					Email:     "kattelsaurav@kami.com",
					Role:      "user",
				},
			},
			want: &lib.ApiResponse{
				Status:  http.StatusBadRequest,
				Message: "password cannot be less the 8 character long",
				Response: lib.Res{
					Error: "invalid phone",
					Data:  nil,
				},
			},
		}, {
			name: "invalid phone",
			args: args{
				payload: &lib.User{
					Id:        "123",
					Username:  "Sam",
					Password:  "1230",
					Phone:     "1234790",
					FirstName: "s",
					LastName:  "23",
					Email:     "kattelsaurav@fmail.com",
					Role:      "user",
				},
			},
			want: &lib.ApiResponse{
				Status:  http.StatusBadRequest,
				Message: "Phone number must be 10 character long",
				Response: lib.Res{
					Error: "invalid phone",
					Data:  nil,
				}},
		}, {
			name: "valid",
			args: args{
				payload: &lib.User{
					Id:        "123",
					Username:  "Sam",
					Password:  "1230873824",
					Phone:     "1234567890",
					FirstName: "s",
					LastName:  "23",
					Email:     "kattelsaurav32@gmail.com",
					Role:      "user",
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.ValidateUserPayload(tt.args.payload); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateUserPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
