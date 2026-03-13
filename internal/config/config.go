package config

import "fmt"

type DBCredentials struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Schema   string
	SSLMode  string
}

func BuildDBUrl(creds DBCredentials) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		creds.Host,
		creds.Port,
		creds.User,
		creds.Password,
		creds.DBName,
		creds.Schema,
		creds.SSLMode,
	)
}
