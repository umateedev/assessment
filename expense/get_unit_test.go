//go:build unit

package expense

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/umateedev/assessment/database"
)

func TestGetExpenseById_ReturnBadRequest_WhenPathMissingId(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues("")

	err := GetExpenseByIdHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestGetExpenseById_ReturnInternalServerError_WhenDbFailed(t *testing.T) {
	e := echo.New()
	expenseId := "1"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues(expenseId)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Open sqlmock error '%s'", err)
	}
	defer db.Close()
	database.Db = db
	mock.ExpectPrepare("SELECT(.*)").
		ExpectQuery().
		WithArgs(expenseId).
		WillReturnError(sqlmock.ErrCancelled)

	err = GetExpenseByIdHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestGetExpenseById_ReturnSuccess(t *testing.T) {
	e := echo.New()
	expenseId := "1"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/expense/:id")
	c.SetParamNames("id")
	c.SetParamValues(expenseId)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Open sqlmock error '%s'", err)
	}
	defer db.Close()

	database.Db = db
	mockExpense := sqlmock.NewRows([]string{"Id", "Title", "Amount", "Note", "Tags"}).
		AddRow("1", "test", 10, "test", pq.Array([]string{"foo", "bar"}))
	mock.ExpectPrepare("SELECT(.*)").
		ExpectQuery().
		WithArgs(expenseId).
		WillReturnRows(mockExpense)

	err = GetExpenseByIdHandler(c)

	expected := "{\"id\":1,\"title\":\"test\",\"amount\":10,\"note\":\"test\",\"tags\":[\"foo\",\"bar\"]}"
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetAllExpense_ReturnInternalServerError_WhenDbFailed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expense", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Open sqlmock error '%s'", err)
	}
	defer db.Close()
	database.Db = db
	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").
		ExpectQuery().
		WillReturnError(sqlmock.ErrCancelled)

	err = GetAllExpenseHandler(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}
}

func TestGetAllExpense_ReturnSuccess(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expense", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Open sqlmock error '%s'", err)
	}
	defer db.Close()

	database.Db = db
	mockExpense := sqlmock.NewRows([]string{"Id", "Title", "Amount", "Note", "Tags"}).
		AddRow("1", "test", 10, "test", pq.Array([]string{"foo", "bar"})).
		AddRow("2", "test2", 10, "test2", pq.Array([]string{"foo2", "bar2"}))
	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").
		ExpectQuery().
		WillReturnRows(mockExpense)

	err = GetAllExpenseHandler(c)

	expected := "[{\"id\":1,\"title\":\"test\",\"amount\":10,\"note\":\"test\",\"tags\":[\"foo\",\"bar\"]},{\"id\":2,\"title\":\"test2\",\"amount\":10,\"note\":\"test2\",\"tags\":[\"foo2\",\"bar2\"]}]"
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
