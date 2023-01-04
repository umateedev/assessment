package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/umateedev/assessment/database"
)

func UpdateExpenseHandler(c echo.Context) error {
	e := Expense{}

	id := c.Param("id")
	if len(id) == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid request, missing param id"})
	}

	err := c.Bind(&e)
	if err != nil {
		log.Printf("Invalid request %s", err.Error())
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid request"})
	}

	stmt, err := database.Db.Prepare("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id")
	if err != nil {
		log.Printf("Prepare statement error", err)
		return c.JSON(http.StatusInternalServerError, "Prepare statement error")
	}

	row := stmt.QueryRow(e.Title, e.Amount, e.Note, pq.Array(&e.Tags), id)
	err = row.Scan(&e.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}
