package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/repository/pgsql_repo"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	pg_user := os.Getenv("POSTGRES_USER")
	if pg_user == "" {
		log.Fatal(".env file error: variable POSTGRES_USER not set")
	}

	pg_pass := os.Getenv("POSTGRES_PASSWORD")
	if pg_pass == "" {
		log.Fatal(".env file error: variable POSTGRES_PASSWORD not set")
	}

	db, err := sqlx.Connect("pgx", "user="+pg_user+" password="+pg_pass+" dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ur := pgsql_repo.NewUserRepository(db)

	admin_name := os.Getenv("ADMIN_USERNAME")
	if admin_name == "" {
		log.Fatal(".env file error: variable ADMIN_USERNAME not set")
	}

	admin, err := ur.FindByName(context.Background(), admin_name)
	if err != nil {
		log.Fatal(err)
	}

	if admin == nil {
		admin_pass := os.Getenv("ADMIN_PASSWORD")
		if admin_pass == "" {
			log.Fatal(".env file error: variable ADMIN_PASSWORD not set")
		}

		admin, err = domain.NewUser(admin_name, admin_pass)
		if err != nil {
			log.Fatal(err)
		}
		err = ur.Save(context.Background(), admin)
		if err != nil {
			log.Fatal(err)
		}
	}
}
