package book // import "book"

import (
	"db"
	"fmt"
	"log"
	"models"
	"net/http"

	"github.com/labstack/echo/v4"
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

// GetBooks - 책정보 취득
func GetBooks(c echo.Context) error {
	data, err := db.SelectData(models.Book{})
	if err != nil {
		log.Fatal("SelectData: ", err)
	}

	// for _, b := range data {
	// 	log.Println(b.IDX, b.Name, b.Author, b.Price)
	// }

	return c.JSON(http.StatusOK, data)
}

// EditBooks - 책정보 수정
func EditBooks(c echo.Context) error {
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
