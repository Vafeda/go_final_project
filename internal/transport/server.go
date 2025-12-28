package transport

import (
	"database/sql"
	"fmt"
	"github.com/Vafeda/go_final_project/internal/transport/handler"
	"net/http"
	"strconv"
	"time"

	"github.com/Vafeda/go_final_project/tests"
)

type Server struct {
	DB   *sql.DB
	HTTP http.Server
}

func Create(db *sql.DB) (*Server, error) {
	addr, err := createAddress()
	fmt.Println(addr)
	if err != nil {
		return nil, err
	}

	server := Server{
		DB: db,
		HTTP: http.Server{
			Addr:         addr,
			Handler:      handler.SetupRoutes(db),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}

	return &server, nil
}

// Стоит сделать логирование
func createAddress() (string, error) {
	addr := ":"
	//if todoPort, exist := os.LookupEnv("TODO_PORT"); exist && len(todoPort) != 0 {
	//	if len(todoPort) != 4 {
	//		// Логирование
	//		fmt.Errorf("TODO_PORT must be exactly 4 characters, got %d characters: %s",
	//			len(todoPort), todoPort)
	//		goto exitCondition
	//	}
	//
	//	if _, err := strconv.Atoi(todoPort); err != nil {
	//		// Логирование
	//		fmt.Errorf("TODO_PORT must contain only digits: %s", todoPort)
	//		goto exitCondition
	//	}
	//
	//	addr += todoPort
	//
	//	return addr, nil
	//}

	//exitCondition:

	addr += strconv.Itoa(tests.Port)

	return addr, nil
}
