//go:build unit

package handler

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateExpensesByID(t *testing.T) {
	t.Run("Should update expenses by id successfully", func(t *testing.T) {
		body := bytes.NewBufferString(`{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`)

		var e Expenses
		res := request(http.MethodPut, uri("expenses", "1"), body)
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
