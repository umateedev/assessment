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

func TestUpdateExpense_ReturnBadRequest_WhenInvalidRequest(t *testing.T) {
	e := echo.New()
	body := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": "food"
	}`
	req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := UpdateExpenseHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateExpense_ReturnBadRequest_WhenPathMissingId(t *testing.T) {
	e := echo.New()
	body := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues("")

	err := UpdateExpenseHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateExpense_ReturnInternalServerError_WhenUpdateFailed(t *testing.T) {
	e := echo.New()
	body := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Open sqlmock error '%s'", err)
	}
	defer db.Close()

	database.Db = db
	mock.ExpectQuery("UPDATE expenses").WillReturnError(sqlmock.ErrCancelled)
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	err = CreateExpenseHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestUpdateExpense_ReturnSuccess(t *testing.T) {
	e := echo.New()
	body := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	updatedExpense := sqlmock.NewRows([]string{"Id"}).AddRow("1")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Open sqlmock error '%s'", err)
	}
	defer db.Close()

	database.Db = db
	mock.ExpectPrepare("UPDATE expenses").ExpectQuery().WillReturnRows(updatedExpense)

	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	err = UpdateExpenseHandler(c)

	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
