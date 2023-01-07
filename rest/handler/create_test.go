//go:build unit

package handler

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (h *handler) TestCreateExpenses(t *testing.T) {
	t.Run("Should create new expenses successfully", func(t *testing.T) {
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
