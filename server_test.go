package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExpenses(t *testing.T) {
	t.Run("Should create new expense successfully", func(t *testing.T) {
		body := bytes.NewBufferString(`{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath", 
			"tags": ["food", "beverage"]
		}`)

		var e Expenses

		res := request(http.MethodPost, uri("expenses"), body)
		err := res.Decode(&e)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, res.StatusCode)
			assert.NotEqual(t, 0, e.ID)
			assert.Equal(t, "strawberry smoothie", e.Title)
			assert.Equal(t, float32(79.00), e.Amount)
			assert.Equal(t, "night market promotion discount 10 bath", e.Note)
			assert.Equal(t, []string{"food", "beverage"}, e.Tags)
		}

	})
}

func TestGetExpensesByID(t *testing.T) {
	t.Run("Should get expenses by id successfully", func(t *testing.T) {
		body := bytes.NewBufferString(`{
			"id": 1
		}`)

		var e Expenses
		res := request(http.MethodGet, uri("expenses", "1"), body)
		err := res.Decode(&e)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, res.StatusCode)
			assert.NotEqual(t, 0, e.ID)
			assert.Equal(t, "strawberry smoothie", e.Title)
			assert.Equal(t, float64(79.00), e.Amount)
			assert.Equal(t, "night market promotion discount 10 bath", e.Note)
			assert.Equal(t, []string{"food", "beverage"}, e.Tags)
		}
	})
}

func TestUpdateExpensesByID(t *testing.T) {
	t.Run("Should update expense by id successfully", func(t *testing.T) {
		body := bytes.NewBufferString(`{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`)

		var e Expenses
		res := request(http.MethodGet, uri("expenses", "1"), body)
		err := res.Decode(&e)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, res.StatusCode)
			assert.Equal(t, 1, e.ID)
			assert.Equal(t, "apple smoothie", e.Title)
			assert.Equal(t, float64(89.00), e.Amount)
			assert.Equal(t, "no discount", e.Note)
			assert.Equal(t, []string{"beverage"}, e.Tags)
		}
	})
}
func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
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

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
