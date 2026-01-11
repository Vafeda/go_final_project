package app

import (
	"github.com/joho/godotenv"

	"github.com/Vafeda/TODO-List/internal/database"
	"github.com/Vafeda/TODO-List/internal/env"
	"github.com/Vafeda/TODO-List/internal/transport"
)

func Run() {
	err := godotenv.Load()
	env.Init()

	db, err := database.Connect("scheduler.db")
	if err != nil {
		panic(err)
	}
	defer database.Close(db)

	srv := transport.NewServer(db)
	if srv == nil {
		return
	}

	srv.Start()
}
