package db

import (
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInitDB(t *testing.T) {
	tests := []struct {
		name    string
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("InitDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// if !reflect.DeepEqual(got, tt.want) {
			if equal := cmp.Equal(got, tt.want); !equal {
				t.Errorf("InitDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTable(t *testing.T) {
	type args struct {
		recreate bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateTable(tt.args.recreate); (err != nil) != tt.wantErr {
				t.Errorf("CreateTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
