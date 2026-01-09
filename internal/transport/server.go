package transport

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Vafeda/TODO-List/internal/transport/handler"
)

const (
	EnvTodoPort = "TODO_PORT"

	DefaultPort = "7540"

	MinPortNum = 1024
	MaxPortNum = 65535
)

type Server struct {
	DB   *sql.DB
	HTTP http.Server
}

func Create(db *sql.DB) (*Server, error) {
	addr, err := createAddress()
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
	todoPort, exist := os.LookupEnv(EnvTodoPort)
	if !exist {
		return "", fmt.Errorf("%s environment variable is not set", EnvTodoPort)
	}

	port, err := strconv.Atoi(todoPort)
	if err != nil {
		return "", fmt.Errorf("%s must contain only digits, got: %s", EnvTodoPort, todoPort)
	}

	if port < MinPortNum || port > MaxPortNum {
		return "", fmt.Errorf("%s value %d is out of valid range (%d-%d)",
			EnvTodoPort, port, MinPortNum, MaxPortNum)
	}

	return todoPort, nil
}
