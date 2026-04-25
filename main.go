package main

import (
	"fmt"

	"github.com/jeagerism/agnos-backend-assignment/database"
	"github.com/jeagerism/agnos-backend-assignment/config"
)

func main() {
	// load config
	cfg := config.LoadConfig()
	// connect to database
	database.ConnectDB(cfg)
	// print server starting message
	fmt.Printf("Server is starting on port %s...\n", cfg.AppPort)
}