package app

import (
	"fmt"
	"github.com/Vafeda/go_final_project/internal/database"
	"github.com/Vafeda/go_final_project/internal/transport"
	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load()

	path := "s"
	//val, exist := os.LookupEnv("TODO_DBFILE")
	//if exist && val != "" {
	//	path = val
	//}

	fmt.Println(path)
	db, err := database.Connect(path)
	if err != nil {
		panic(err)
	}
	defer database.Disconnect(db)

	srv, err := transport.Create(db)
	if err != nil {
		panic(err)
	}

	if err = srv.HTTP.ListenAndServe(); err != nil {
		fmt.Println(err)
		return
	}
}
