package test

import (
	"bytes"
	"io"
	"net/http"
	"sauravkattel/agri/src/lib"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProductPayload(t *testing.T) {
	type args struct {
		r *http.Request
	}
	type T any

	tests := []struct {
		name    string
		args    args
		want    *lib.ProductPayload
		wantErr bool
	}{
		{
			name: "valid test",
			args: args{
				r: &http.Request{
					Body: io.NopCloser(bytes.NewBufferString(`{			
						"product": {
						"name": "johndoe",
						"description": "jejec"
						},
						"attrib": {
						"price": 1.2,
						"quantity": 200,
						"status": "1"
						}
					}`)),
				},
			},

			want: &lib.ProductPayload{
				Product: lib.Product{
					Name:        "johndoe",
					Description: "jejec",
				},
				Attrib: lib.Attrib{
					Status:   "1",
					Quantity: 200,
					Price:    1.2,
				},
			},
			wantErr: false,
		}, {
			name: "invalid test",
			args: args{
				r: &http.Request{
					Body: io.NopCloser(bytes.NewBufferString(``)),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lib.ParseJson[lib.ProductPayload](tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want, "the retrun value was nil")
		})
	}
}
func TestParseJsonUserPayload(t *testing.T) {
	type args struct {
		r *http.Request
	}
	type T any

	tests := []struct {
		name    string
		args    args
		want    *lib.UserPayload
		wantErr bool
	}{
		{
			name: "valid test",
			args: args{
				r: &http.Request{
					Body: io.NopCloser(bytes.NewBufferString(`{
						"username": "johndoe",
						"email": "john.doe@example.com",
						"password": "password123",
						"phone": "1234567890",
						"account_type": "user",
						"address": "123 Main St"
					}`)),
				},
			},
			want: &lib.UserPayload{
				Username:    "johndoe",
				Email:       "john.doe@example.com",
				Password:    "password123",
				Phone:       "1234567890",
				AccountType: "user",
				Address:     "123 Main St",
			},
			wantErr: false,
		}, {
			name: "invalid test",
			args: args{
				r: &http.Request{
					Body: io.NopCloser(bytes.NewBufferString(``)),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lib.ParseJson[lib.UserPayload](tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want, "the retrun value was nil")
		})
	}
}

func TestPasreString(t *testing.T) {
	type T any
	type args struct {
		data *lib.UserPayload
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "valid data",
			args: args{
				data: &lib.UserPayload{
					Username:    "johndoe",
					Email:       "john.doe@example.com",
					Password:    "password123",
					Phone:       "1234567890",
					AccountType: "user",
					Address:     "123 Main St",
				},
			},
			want:    []byte(`{"username":"johndoe","email":"john.doe@example.com","password":"password123","phone":"1234567890","account_type":"user","address":"123 Main St"}`),
			wantErr: false,
		},
		{
			name: "empty data",
			args: args{
				data: &lib.UserPayload{},
			},
			want:    []byte(`{"username":"","email":"","password":"","phone":"","account_type":"","address":""}`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lib.PasreString(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasreString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.JSONEq(t, string(tt.want), string(got), "ParseString() got = %v, want %v", string(got), string(tt.want))
		})
	}
}
