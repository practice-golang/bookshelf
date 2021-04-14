package db

import (
	"reflect"
	"testing"

	"github.com/doug-martin/goqu/v9"
)

func TestPrepareWhere(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want goqu.Ex
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrepareWhere(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrepareWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckValidAndPrepareWhere(t *testing.T) {
	type args struct {
		book interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    goqu.Ex
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckValidAndPrepareWhere(tt.args.book)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckValidAndPrepareWhere() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckValidAndPrepareWhere() = %v, want %v", got, tt.want)
			}
		})
	}
}
