package app

import (
	"github.com/Vafeda/go_final_project/internal/database"
	"github.com/Vafeda/go_final_project/internal/transport"
	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load()

	db, err := database.Connect("scheduler.db")
	if err != nil {
		panic(err)
	}
	defer database.Disconnect(db)

	srv, err := transport.Create(db)
	if err != nil {
		panic(err)
	}

	if err = srv.HTTP.ListenAndServe(); err != nil {
		return
	}
}
