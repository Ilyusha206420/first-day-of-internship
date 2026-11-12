package main

import "github.com/Ilyusha206420/first-day-of-internship/backend/internal"

func main() {
	db, err := internal.ConnectDB()
	if err != nil {
		panic(err)
	}
	uc, err := internal.NewUseCase(db)
	if err != nil {
		panic(err)
	}
}
