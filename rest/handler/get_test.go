//go:build unit

package handler

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestGetExpenses(t *testing.T) {
	t.Run("Should get expenses successfully", func(t *testing.T) {
		body := bytes.NewBufferString(`{}`)

		var e []Expenses
		res := request(http.MethodPut, uri("expenses"), body)
		err := res.Decode(&e)

		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusOK, res.StatusCode)
		assert.Greater(t, len(e), 0)
	})
}
