package app

import (
	"fmt"
	"github.com/joho/godotenv"

	"github.com/Vafeda/TODO-List/internal/database"
	"github.com/Vafeda/TODO-List/internal/transport"
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
	
	fmt.Printf("Serving on port %s\n", srv.HTTP.Addr)
	if err = srv.HTTP.ListenAndServe(); err != nil {
		return
	}
}
