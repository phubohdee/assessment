package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetExpensesByID(c echo.Context) error {
	id := c.Param("id")
	stmt, err := h.db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses statment:" + err.Error()})
	}

	row := stmt.QueryRow(id)
	e := Expenses{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, &e.Tags)
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "id not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expenses:" + err.Error()})
	}
}

func (h *handler) getExpenses(c echo.Context) error {
	stmt, err := h.db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "can't prepare query all expenses statement:" + err.Error()})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "can't query all expenses:" + err.Error()})
	}

	expenses := []Expenses{}

	for rows.Next() {
		var e Expenses

		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, &e.Tags)

		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: "can't scan expenses:" + err.Error()})
		}

		expenses = append(expenses, e)
	}

	return c.JSON(http.StatusOK, expenses)
}
