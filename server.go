package main

import (
	"database/sql"
	"log"
)

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

	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))
}
