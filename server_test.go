package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
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
	t.Run("Should get expense by id successfully", func(t *testing.T) {
		c := seedUser(t)

		var latest Expenses
		res := request(http.MethodGet, uri("users", strconv.Itoa(c.ID)), nil)
		err := res.Decode(&latest)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.NotEqual(t, c.ID, latest.ID)
		assert.NotEmpty(t, latest.Title)
		assert.NotEmpty(t, latest.Amount)
		assert.NotEmpty(t, latest.Note)
		assert.NotEmpty(t, latest.Tags)
	})
}

func seedUser(t *testing.T) Expenses {
	var e Expenses
	body := bytes.NewBufferString(`{
		"name": "Phubohdee",
		"age": 22
	}`)

	err := request(http.MethodPost, uri("expenses"), body).Decode(&e)
	if err != nil {
		t.Fatal("can't create uomer", err)
	}

	return e
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
