package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/lib/pq"
)

type Expenses struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   int      `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

func CreateExpenses(c echo.Context) error {
	var e Expenses
	err := c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags))
	err = row.Scan(&e.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}

func GetExpensesByID(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
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

func UpdateExpensesByID(c echo.Context) error {
	var e Expenses
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	id := c.Param("id")
	stmt, err := db.Prepare("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id")
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

func getExpenses(c echo.Context) error {
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
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

var db *sql.DB

func main() {
	// url := os.Getenv("DATABASE_URL")
	url := "postgres://ulbqilyz:40D21uWaCiKxFMYSCA7-KPaopqdhLAJZ@john.db.elephantsql.com/ulbqilyz"

	var err error
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g := e.Group("expenses")
	g.POST("", CreateExpenses)
	g.GET("/:id", GetExpensesByID)
	g.PUT("/:id", UpdateExpensesByID)
	g.GET("", getExpenses)

	log.Fatal(e.Start(os.Getenv("PORT")))
}
