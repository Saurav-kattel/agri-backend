package test

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"sauravkattel/agri/src/database"
)

func TestConnect(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	type args struct {
		dbUsername string
		dbName     string
		password   string
	}
	tests := []struct {
		name    string
		args    args
		want    *database.Database
		wantErr bool
	}{
		{
			name: "connecion test",
			args: args{
				dbUsername: "postgres",
				dbName:     "agri",
				password:   "saurav",
			},
			want: &database.Database{
				DB: sqlxDB,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := database.Connect(tt.args.dbUsername, tt.args.dbName, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Connect() = %v, want %v", got, tt.want)
			}
		})
	}
}
