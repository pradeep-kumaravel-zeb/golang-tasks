package main

import (
	"fmt"
	"os"
	"student-enrollment/internal/config"
	"student-enrollment/pkg/server"
)

func main() {
	creds, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	app, err := server.InitializeApp(creds)
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server starting on :8080")
	if err := app.Start(); err != nil {
		fmt.Printf("Server failed: %v\n", err)
		os.Exit(1)
	}
}
