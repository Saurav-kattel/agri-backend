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

func TestCreateUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	/*
		payload.FirstName,
			payload.LastName,
			payload.Email,
			hash,
			payload.Phone,
			payload.Username,
			role,
	*/

	mock.ExpectExec("INSERT INTO users").WithArgs("Saurav", "Kattel", "sauravkattel@32gmail.com", "23132jdsadas", "12341213131", "asurab", "user").WillReturnResult(sqlmock.NewResult(1, 1))
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
					FirstName: "Saurav",
					LastName:  "Kattel",
					Phone:     "12341213131",
					Email:     "sauravkattel@32gmail.com",
					Username:  "asurab",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		err := users.CreateUser(sqlxDB, tt.args.payload, tt.args.hash, tt.args.role)
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
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
		Password:  "hashedpassword",
		Role:      "user",
		Username:  "johndoe",
		CreatedAt: &createdAt,
	}

	// Define columns in the same order as the query
	rows := mock.NewRows([]string{
		"id", "first_name", "last_name", "email", "phone", "password", "role", "created_at", "username",
	}).AddRow(
		expectedUser.Id, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email,
		expectedUser.Phone, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.Username,
	)

	query := `SELECT \* FROM users WHERE id = \$1`
	mock.ExpectQuery(query).
		WithArgs("1").
		WillReturnRows(rows)

	user, err := users.GetUsersByUserId(db, "1")
	if err != nil {
		t.Fatalf("Error fetching user: %v", err)
	}

	assert.NotNil(t, user, "Expected user to be found")
	assert.Equal(t, expectedUser.Id, user.Id, "Unexpected user ID")
	assert.Equal(t, expectedUser.FirstName, user.FirstName, "Unexpected first name")
	assert.Equal(t, expectedUser.LastName, user.LastName, "Unexpected last name")
	assert.Equal(t, expectedUser.Email, user.Email, "Unexpected email")
	assert.Equal(t, expectedUser.Phone, user.Phone, "Unexpected phone")
	assert.Equal(t, expectedUser.Password, user.Password, "Unexpected password")
	assert.Equal(t, expectedUser.Role, user.Role, "Unexpected role")
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
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
		Password:  "hashedpassword",
		Role:      "user",
		Username:  "johndoe",
		CreatedAt: &createdAt,
	}

	// Define columns in the same order as the query
	rows := mock.NewRows([]string{
		"id", "first_name", "last_name", "email", "phone", "password", "role", "created_at", "username",
	}).AddRow(
		expectedUser.Id, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email,
		expectedUser.Phone, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.Username,
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
	assert.Equal(t, expectedUser.FirstName, user.FirstName, "Unexpected first name")
	assert.Equal(t, expectedUser.LastName, user.LastName, "Unexpected last name")
	assert.Equal(t, expectedUser.Email, user.Email, "Unexpected email")
	assert.Equal(t, expectedUser.Phone, user.Phone, "Unexpected phone")
	assert.Equal(t, expectedUser.Password, user.Password, "Unexpected password")
	assert.Equal(t, expectedUser.Role, user.Role, "Unexpected role")
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
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
		Password:  "hashedpassword",
		Role:      "user",
		Username:  "johndoe",
		CreatedAt: &createdAt,
	}

	// Define columns in the same order as the query
	rows := mock.NewRows([]string{
		"id", "first_name", "last_name", "email", "phone", "password", "role", "created_at", "username",
	}).AddRow(
		expectedUser.Id, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email,
		expectedUser.Phone, expectedUser.Password, expectedUser.Role, expectedUser.CreatedAt, expectedUser.Username,
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
	assert.Equal(t, expectedUser.FirstName, user.FirstName, "Unexpected first name")
	assert.Equal(t, expectedUser.LastName, user.LastName, "Unexpected last name")
	assert.Equal(t, expectedUser.Email, user.Email, "Unexpected email")
	assert.Equal(t, expectedUser.Phone, user.Phone, "Unexpected phone")
	assert.Equal(t, expectedUser.Password, user.Password, "Unexpected password")
	assert.Equal(t, expectedUser.Role, user.Role, "Unexpected role")
	assert.Equal(t, expectedUser.Username, user.Username, "Unexpected username")

	// Ensure all expectations are fulfilled
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "Unfulfilled expectations")
}
