package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phubohdee/assessment/rest/handler"
)

var h *handler.Handler

func main() {
	url := os.Getenv("DATABASE_URL")

	var err error
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g := e.Group("expenses")
	g.POST("", h.CreateExpenses)
	g.GET("/:id", h.GetExpensesByID)
	g.PUT("/:id", h.UpdateExpensesByID)
	g.GET("", h.getExpenses)

	log.Fatal(e.Start(os.Getenv("PORT")))
}
