package book // import "book"

import (
	"db"
	"fmt"
	"log"
	"math"
	"models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

// AddBooks - Insert book(s) info
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

// GetBook - Get a book info
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

// GetBooks - Get all(but limit 10 by db.SelectData) books info
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

// SearchBooks - Search book(s) info or paging
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

// EditBook - Edit book info
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

// DeleteBook - Delete an item of books info
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

// GetTotalPage - Get total page
func GetTotalPage(c echo.Context) error {
	var search models.BookSearch

	if err := c.Bind(&search); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	data, err := db.SelectCount(search)
	if err != nil {
		log.Fatal("SelectCount: ", err)
	}

	countPerPage := uint(1)
	if search.Options.Count.Valid {
		countPerPage = uint(search.Options.Count.Int64)
	}

	pages := uint(math.Ceil(float64(data) / float64(countPerPage)))

	result := map[string]uint{"total-page": pages}

	return c.JSON(http.StatusOK, result)
}
