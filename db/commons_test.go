package db

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestInsertData(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    sql.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectData(t *testing.T) {
	type args struct {
		search interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectData(tt.args.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateData(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    sql.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteData(t *testing.T) {
	type args struct {
		target string
		value  string
	}
	tests := []struct {
		name    string
		args    args
		want    sql.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteData(tt.args.target, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectCount(t *testing.T) {
	tests := []struct {
		name    string
		want    uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectCount()
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SelectCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
