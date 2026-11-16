package main

import (
	"context"
	"log"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/repository/pgsql_repo"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("pgx", "user=postgres password=password dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ur := pgsql_repo.NewUserRepository(db)
	admin, err := ur.FindByName(context.Background(), "admin")
	if err != nil {
		log.Fatal(err)
	}
	if admin == nil {
		admin, err = domain.NewUser("admin", "admin123")
		if err != nil {
			log.Fatal(err)
		}
		err = ur.Save(context.Background(), admin)
		if err != nil {
			log.Fatal(err)
		}
	}
}
