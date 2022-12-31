package expense

func main() {

}

func (h *handler) CreateExpense(c echo.Context) error {
	e := Expense{}
	err := c.Bind(e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "Invalid request"})
	}

	row := database.db.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags))
	err := row.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
