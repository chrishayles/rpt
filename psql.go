package rpt

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresClient struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Client   sql.DB
}

func NewPostgresClient() *PostgresClient {
	return &PostgresClient{}
}

func (psql *PostgresClient) Connect() error    { return nil }
func (psql *PostgresClient) Disconnect() error { return nil }
func (psql *PostgresClient) Seed(d DataSet) error {
	fmt.Println("Seeding...")
	fmt.Println(string(d.ToJson()))
	return nil
}
