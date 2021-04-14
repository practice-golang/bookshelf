package book

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAddBooks(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := AddBooks(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("AddBooks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBook(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := GetBook(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBooks(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := GetBooks(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetBooks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchBooks(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := SearchBooks(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SearchBooks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEditBook(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := EditBook(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("EditBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := DeleteBook(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTotalPage(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := GetTotalPage(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetTotalPage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
