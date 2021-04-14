package main

import (
	"testing"

	"github.com/labstack/echo/v4"
	_ "modernc.org/sqlite"

	"github.com/google/go-cmp/cmp"
)

func Test_setupDB(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setupDB(); (err != nil) != tt.wantErr {
				t.Errorf("setupDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_dumpHandler(t *testing.T) {
	type args struct {
		c       echo.Context
		reqBody []byte
		resBody []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dumpHandler(tt.args.c, tt.args.reqBody, tt.args.resBody)
		})
	}
}

func Test_setupServer(t *testing.T) {
	tests := []struct {
		name string
		want *echo.Echo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := setupServer(); !reflect.DeepEqual(got, tt.want) {
			got := setupServer()
			if equal := cmp.Equal(got, tt.want); !equal {
				t.Errorf("setupServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
