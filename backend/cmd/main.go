package main

import (
	"log"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/repository/pgsql_repo"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ur := pgsql_repo.NewUserRepository(db)
}
