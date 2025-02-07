package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/umateedev/assessment/database"
)

func GetExpenseByIdHandler(c echo.Context) error {
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

	}
}

func GetAllExpenseHandler(c echo.Context) error {

	stmt, err := database.Db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		log.Printf("Prepare statement error", err)
		return c.JSON(http.StatusInternalServerError, "Prepare statement error")
	}

	expenses := []Expense{}
	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	for rows.Next() {
		e := Expense{}
		err := rows.Scan(&e.Id, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})

		}
		expenses = append(expenses, e)
	}

	if len(expenses) == 0 {
		return c.JSON(http.StatusNotFound, Error{Message: "expense not found"})
	}

	return c.JSON(http.StatusOK, expenses)
}
