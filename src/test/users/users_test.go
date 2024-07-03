package test

import (
	"net/http"
	"reflect"
	"sauravkattel/agri/src/lib"
	"sauravkattel/agri/src/users"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestValidateUserPayload(t *testing.T) {
	type args struct {
		payload *lib.UserPayload
	}
	tests := []struct {
		name string
		args args
		want *lib.ApiResponse
	}{
		{
			name: "invalid email",
			args: args{
				payload: &lib.UserPayload{
					Username:    "Sam",
					Password:    "1234567890",
					Email:       "kattelsaurav.com",
					Phone:       "1234567890",
					Address:     "ilam",
					AccountType: "buyer",
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
				payload: &lib.UserPayload{
					Username:    "Sam",
					Password:    "1290",
					Email:       "kattelsaurav32@gmail.com",
					Phone:       "1234567890",
					Address:     "ilam",
					AccountType: "buyer",
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
				payload: &lib.UserPayload{
					Username:    "Sam",
					Password:    "1234567890",
					Email:       "kattelsaurav32@gmail.com",
					Phone:       "12789",
					Address:     "ilam",
					AccountType: "buyer",
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
				payload: &lib.UserPayload{
					Username:    "Sam",
					Password:    "1234567890",
					Email:       "kattelsaurav@32gmail.com",
					Phone:       "1234567890",
					Address:     "ilam",
					AccountType: "buyer",
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

func TestCreateUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectQuery("INSERT INTO users").WithArgs("sauravkattel@32gmail.com", "23132jdsadas", "12341213131", "asurab", "buyer", "ilam").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	type args struct {
		hash    string
		role    string
		payload *lib.UserPayload
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create user tesst",
			args: args{
				hash: "23132jdsadas",
				role: "user",
				payload: &lib.UserPayload{
					Phone:       "12341213131",
					Email:       "sauravkattel@32gmail.com",
					Username:    "asurab",
					AccountType: "buyer",
					Address:     "ilam",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		_, err := users.CreateUser(sqlxDB, tt.args.payload, tt.args.hash)
		if (err != nil) != tt.wantErr {
			t.Errorf("CreateUsers error occured %+v want err %+v", err, tt.wantErr)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expections were not fullfillled %+v", err)
	}
}

func TestGetUsersByUserId(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	createdAt := "random-date-2342"
	expectedUser := lib.User{
		Id:          "1",
		Email:       "john.doe@example.com",
		Phone:       "1234567890",
		Password:    "hashedpassword",
		AccountType: "buyer",
		Address:     "ilam",
		Username:    "johndoe",
		CreatedAt:   &createdAt,
	}

	// Define columns in the same order as the query
	rows := mock.NewRows([]string{
		"id", "email", "phone", "password", "created_at", "username", "account_type", "address",
	}).AddRow(
		expectedUser.Id, expectedUser.Email,
		expectedUser.Phone, expectedUser.Password, expectedUser.CreatedAt, expectedUser.Username, expectedUser.AccountType, expectedUser.Address,
	)

	query := `SELECT \* FROM users WHERE email = \$1`
	mock.ExpectQuery(query).
		WithArgs("john.doe@example.com").
		WillReturnRows(rows)

	user, err := users.GetUsersByEmail(db, "john.doe@example.com")
	if err != nil {
		t.Fatalf("Error fetching user: %v", err)
	}

	assert.NotNil(t, user, "Expected user to be found")
	assert.Equal(t, expectedUser.Id, user.Id, "Unexpected user ID")
	assert.Equal(t, expectedUser.AccountType, user.AccountType, "Unexpcted Account type")
	assert.Equal(t, expectedUser.Address, user.Address, "unexpcted address")
	assert.Equal(t, expectedUser.Email, user.Email, "Unexpected email")
	assert.Equal(t, expectedUser.Phone, user.Phone, "Unexpected phone")
	assert.Equal(t, expectedUser.Password, user.Password, "Unexpected password")
	assert.Equal(t, expectedUser.Username, user.Username, "Unexpected username")

	// Ensure all expectations are fulfilled
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Unfulfilled expectations")
}

func TestGetUsersByEmail(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	createdAt := "random-date-2342"
	expectedUser := lib.User{
		Id:          "1",
		Email:       "john.doe@example.com",
		Phone:       "1234567890",
		Password:    "hashedpassword",
		AccountType: "buyer",
		Address:     "ilam",
		Username:    "johndoe",
		CreatedAt:   &createdAt,
	}

	// Define columns in the same order as the query
	rows := mock.NewRows([]string{
		"id", "email", "phone", "password", "created_at", "username", "account_type", "address",
	}).AddRow(
		expectedUser.Id, expectedUser.Email,
		expectedUser.Phone, expectedUser.Password, expectedUser.CreatedAt, expectedUser.Username, expectedUser.AccountType, expectedUser.Address,
	)

	query := `SELECT \* FROM users WHERE email = \$1`
	mock.ExpectQuery(query).
		WithArgs("john.doe@example.com").
		WillReturnRows(rows)

	user, err := users.GetUsersByEmail(db, "john.doe@example.com")
	if err != nil {
		t.Fatalf("Error fetching user: %v", err)
	}

	assert.NotNil(t, user, "Expected user to be found")
	assert.Equal(t, expectedUser.Id, user.Id, "Unexpected user ID")
	assert.Equal(t, expectedUser.AccountType, user.AccountType, "Unexpcted Account type")
	assert.Equal(t, expectedUser.Address, user.Address, "unexpcted address")
	assert.Equal(t, expectedUser.Email, user.Email, "Unexpected email")
	assert.Equal(t, expectedUser.Phone, user.Phone, "Unexpected phone")
	assert.Equal(t, expectedUser.Password, user.Password, "Unexpected password")
	assert.Equal(t, expectedUser.Username, user.Username, "Unexpected username")

	// Ensure all expectations are fulfilled
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Unfulfilled expectations")
}

func TestGetUsersByUserName(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	createdAt := "random-date-2342"
	expectedUser := lib.User{
		Id:          "1",
		Email:       "john.doe@example.com",
		Phone:       "1234567890",
		Password:    "hashedpassword",
		AccountType: "buyer",
		Address:     "ilam",
		Username:    "johndoe",
		CreatedAt:   &createdAt,
	}

	// Define columns in the same order as the query
	rows := mock.NewRows([]string{
		"id", "email", "phone", "password", "created_at", "username", "account_type", "address",
	}).AddRow(
		expectedUser.Id, expectedUser.Email,
		expectedUser.Phone, expectedUser.Password, expectedUser.CreatedAt, expectedUser.Username, expectedUser.AccountType, expectedUser.Address,
	)

	query := `SELECT \* FROM users WHERE username = \$1`
	mock.ExpectQuery(query).
		WithArgs("johndoe").
		WillReturnRows(rows)

	user, err := users.GetUsersByUserName(db, "johndoe")
	if err != nil {
		t.Fatalf("Error fetching user: %v", err)
	}

	assert.NotNil(t, user, "Expected user to be found")
	assert.Equal(t, expectedUser.Id, user.Id, "Unexpected user ID")
	assert.Equal(t, expectedUser.AccountType, user.AccountType, "Unexpcted Account type")
	assert.Equal(t, expectedUser.Address, user.Address, "unexpcted address")
	assert.Equal(t, expectedUser.Email, user.Email, "Unexpected email")
	assert.Equal(t, expectedUser.Phone, user.Phone, "Unexpected phone")
	assert.Equal(t, expectedUser.Password, user.Password, "Unexpected password")
	assert.Equal(t, expectedUser.Username, user.Username, "Unexpected username")

	// Ensure all expectations are fulfilled
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Unfulfilled expectations")
}
