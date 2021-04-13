package models // import "models"

import (
	"gopkg.in/guregu/null.v4"
)

type Book struct {
	Idx    null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	Name   null.String `json:"name" db:"NAME"`
	Price  null.Float  `json:"price" db:"PRICE"`
	Author null.String `json:"author" db:"AUTHOR"`
	ISBN   null.String `json:"isbn" db:"ISBN"`
}
