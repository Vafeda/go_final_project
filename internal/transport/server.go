package transport

import (
	"database/sql"
	"fmt"
	"github.com/Vafeda/go_final_project/internal/transport/handler"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
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
	todoPort, exist := os.LookupEnv("TODO_PORT")
	if !exist {
		return "", fmt.Errorf("TODO_PORT environment variable is not set")
	}

	port, err := strconv.Atoi(todoPort)
	if err != nil {
		return "", fmt.Errorf("TODO_PORT must contain only digits, got: %s", todoPort)
	}

	if port < MinPortNum || port > MaxPortNum {
		return "", fmt.Errorf("TODO_PORT value %d is out of valid range (%d-%d)",
			port, MinPortNum, MaxPortNum)
	}

	return todoPort, nil
}
