//go:build unit

package expense

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/umateedev/assessment/database"
)

func TestCreateExpense_ReturnBadRequest_WhenInvalidRequest(t *testing.T) {
	e := echo.New()
	body := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": "food"
	}`
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := CreateExpenseHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateExpense_ReturnInternalServerError_WhenInsertFailed(t *testing.T) {
	e := echo.New()
	body := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Open sqlmock error '%s'", err)
	}
	defer db.Close()

	database.Db = db
	mock.ExpectQuery("INSERT INTO expenses").WillReturnError(sqlmock.ErrCancelled)
	c := e.NewContext(req, rec)

	err = CreateExpenseHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestCreateExpense_ReturnSuccess(t *testing.T) {
	e := echo.New()
	body := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newExpense := sqlmock.NewRows([]string{"Id"}).AddRow("1")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Open sqlmock error '%s'", err)
	}
	defer db.Close()

	database.Db = db
	mock.ExpectQuery("INSERT INTO expenses").WillReturnRows(newExpense)
	c := e.NewContext(req, rec)
	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	err = CreateExpenseHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
