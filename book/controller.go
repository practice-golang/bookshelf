package book // import "book"

import (
	"db"
	"fmt"
	"log"
	"models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

// AddBooks - 책정보 입력
func AddBooks(c echo.Context) error {
	var books []models.Book

	if err := c.Bind(&books); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	sqlResult, err := db.InsertData(books)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// GetBook - 책정보 한개 취득
func GetBook(c echo.Context) error {
	idx := c.Param("idx")
	// var data interface{}
	var data []models.Book
	var err error
	dataINTF, err := db.SelectData(models.Book{Idx: null.NewString(idx, true)})
	if err != nil {
		log.Fatal("SelectData: ", err)
	}
	if dataINTF != nil {
		data = dataINTF.([]models.Book)
	} else {
		log.Println("Control WTF null")
	}
	// data, err = db.SelectData(models.Book{})
	// if err != nil {
	// 	log.Fatal("SelectData: ", err)
	// }

	return c.JSON(http.StatusOK, data)
}

// GetBooks - 책정보 취득 검색
func GetBooks(c echo.Context) error {
	// var data interface{}
	var data []models.Book
	var err error
	dataINTF, err := db.SelectData(models.Book{})
	if err != nil {
		log.Fatal("SelectData: ", err)
	}
	if dataINTF != nil {
		data = dataINTF.([]models.Book)
	} else {
		log.Println("Control WTF null")
	}

	return c.JSON(http.StatusOK, data)
}

// SearchBooks - 책정보 검색 또는 페이징
func SearchBooks(c echo.Context) error {
	var data []models.Book
	var search models.BookSearch
	var err error

	if err := c.Bind(&search); err != nil {
		log.Fatal("Search data: ", err)
	}

	dataINTF, err := db.SelectData(search)
	if err != nil {
		log.Fatal("SelectData: ", err)
	}
	if dataINTF != nil {
		data = dataINTF.([]models.Book)
	} else {
		log.Println("Control WTF null")
	}

	return c.JSON(http.StatusOK, data)
}

// EditBook - 책정보 수정
func EditBook(c echo.Context) error {
	var book models.Book

	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	sqlResult, err := db.UpdateData(book)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteBook - 책 1개 삭제
func DeleteBook(c echo.Context) error {
	idx := c.Param("idx")

	sqlResult, err := db.DeleteData("idx", idx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}
