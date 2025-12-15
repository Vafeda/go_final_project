package app

import (
	"fmt"
	"github.com/Vafeda/go_final_project/internal/transport"
)

func Run() {
	srv, err := transport.Create()
	if err != nil {
		panic(err)
	}

	if err = srv.HTTP.ListenAndServe(); err != nil {
		fmt.Println(err)
		return
	}
}
