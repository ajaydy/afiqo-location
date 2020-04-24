package helpers

import (
	"database/sql"
	"fmt"
)

type DBOptions struct {
	Host        string
	Port        int
	Username    string
	Password    string
	DBName      string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
	SSLMode     string
}

func InitDB(options DBOptions) (*sql.DB, error) {
	dbConfig := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		options.Username,
		options.Password,
		options.Host,
		options.Port,
		options.DBName,
		options.SSLMode)
	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
