package transport

import (
	"database/sql"
	"fmt"
	"github.com/Vafeda/TODO-List/internal/env"
	"net/http"
	"strconv"
	"time"

	"github.com/Vafeda/TODO-List/internal/transport/handler"
)

const (
	EnvTodoPort = "TODO_PORT"

	DefaultPort = "7540"
)

type Server struct {
	db   *sql.DB
	http http.Server
}

func NewServer(db *sql.DB) *Server {
	addr, err := createAddress()
	if err != nil {
		return nil
	}

	server := Server{
		db: db,
		http: http.Server{
			Addr:         addr,
			Handler:      handler.SetupRoutes(db),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}

	return &server
}

func (s *Server) Start() {
	fmt.Printf("Serving on port %s\n", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil {
		return
	}
}

func createAddress() (string, error) {
	var addr string

	port, err := getAddrFromEnvPort()
	if err != nil {
		addr = ":" + DefaultPort
	} else {
		addr = ":" + port
	}

	return addr, nil
}

func getAddrFromEnvPort() (string, error) {
	todoPort, exist := env.Dict[EnvTodoPort]
	if !exist {
		return "", fmt.Errorf("%s environment variable is not set", EnvTodoPort)
	}

	_, err := strconv.Atoi(todoPort)
	if err != nil {
		return "", fmt.Errorf("%s must contain only digits, got: %s", EnvTodoPort, todoPort)
	}

	return todoPort, nil
}
