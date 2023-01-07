package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpensesByID(c echo.Context) error {
	var e Expenses
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	id := c.Param("id")
	stmt, err := h.db.Prepare("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses statment:" + err.Error()})
	}

	row := stmt.QueryRow(e.Title, e.Amount, e.Note, pq.Array(&e.Tags), id)
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, &e.Tags)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't update expenses: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
