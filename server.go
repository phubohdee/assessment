package main

import (
	"database/sql"
	"log"
	"net/http"

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

	e.POST("/expenses", CreateExpenses)

	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))
}
