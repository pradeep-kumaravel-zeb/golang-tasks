package main

import (
	"flag"
	"fmt"
	"os"
	"student-enrollment/pkg/server"
)

func main() {
	host := flag.String("host", "localhost", "Database host")
	port := flag.String("port", "5432", "Database port")
	user := flag.String("user", "postgres", "Database user")
	password := flag.String("password", "root", "Database password")
	dbname := flag.String("dbname", "postgres", "Database name")
	schema := flag.String("schema", "public", "Database schema")
	flag.Parse()

	app, err := server.InitializeApp(*host, *port, *user, *password, *dbname, *schema)
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
