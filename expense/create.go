package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/umateedev/assessment/database"
)

func CreateExpenseHandler(c echo.Context) error {
	e := Expense{}
	err := c.Bind(&e)
	if err != nil {
		log.Printf("Invalid request %s", err.Error())
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid request"})
	}

	row := database.Db.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags))
	err = row.Scan(&e.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
