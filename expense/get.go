package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/umateedev/assessment/database"
)

func GetExpenseHandler(c echo.Context) error {
	id := c.Param("id")
	if len(id) == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid request, missing param id"})
	}

	stmt, err := database.Db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id=$1")
	if err != nil {
		log.Printf("Prepare statement error", err)
		return c.JSON(http.StatusInternalServerError, "Prepare statement error")
	}

	e := Expense{}
	row := stmt.QueryRow(id)
	err = row.Scan(&e.Id, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Error{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Error{Message: "can't scan expense:" + err.Error()})

		return c.JSON(http.StatusOK, e)
	}
}
