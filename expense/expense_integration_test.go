//go:build integration

package expense

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExponse(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "test title",
		"amount": 79,
		"note": "test note", 
		"tags": ["foo", "bar"]
	}`)
	var e Expense

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.Id)
	assert.Equal(t, "test title", e.Title)
	assert.Equal(t, 79.0, e.Amount)
	assert.Equal(t, "test note", e.Note)
	assert.Equal(t, []string{"foo", "bar"}, e.Tags)
}

func TestGetExponse(t *testing.T) {
	e := seedExpense(t)

	var latest Expense

	res := request(http.MethodGet, uri("expenses", strconv.Itoa(e.Id)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.Id, latest.Id)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)
}

func seedExpense(t *testing.T) Expense {
	var c Expense
	body := bytes.NewBufferString(`{
		"title": "seed expense",
		"amount": 79.0,
		"note": "test note", 
		"tags": ["food", "beverage"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return c
}

func uri(paths ...string) string {
	host := "http://localhost" + os.Getenv("PORT")
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}