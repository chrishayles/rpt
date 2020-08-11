package rpt

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

/*

https://pkg.go.dev/github.com/lib/pq?tab=doc

For compatibility with libpq, the following special connection parameters are supported:

* dbname - The name of the database to connect to
* user - The user to sign in as
* password - The user's password
* host - The host to connect to. Values that start with / are for unix
  domain sockets. (default is localhost)
* port - The port to bind to. (default is 5432)
* sslmode - Whether or not to use SSL (default is require, this is not
  the default for libpq)
* fallback_application_name - An application_name to fall back to if one isn't provided.
* connect_timeout - Maximum wait for connection, in seconds. Zero or
  not specified means wait indefinitely.
* sslcert - Cert file location. The file must contain PEM encoded data.
* sslkey - Key file location. The file must contain PEM encoded data.
* sslrootcert - The location of the root certificate file. The file
  must contain PEM encoded data.
* spn - Configures GSS (Kerberos) SPN.
* service - GSS (Kerberos) service name to use when constructing the SPN (default is `postgres`).
Valid values for sslmode are:

* disable - No SSL
* require - Always SSL (skip verification)
* verify-ca - Always SSL (verify that the certificate presented by the
  server was signed by a trusted CA)
* verify-full - Always SSL (verify that the certification presented by
  the server was signed by a trusted CA and the server host name
  matches the one in the certificate)

*/

type PostgresClient struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Client   *sql.DB
	Logger   *Logger
}

type PostgresDatabase struct {
	Datname       string
	Datdba        int
	Encoding      int
	Datcollate    string
	Datctype      string
	Datistemplate bool
	Datallowconn  bool
	Datconnlimit  int
	Datlastsysoid int
	Datfrozenxid  int
	Datminmxid    int
	Dattablespace int
	Datacl        interface{}
}

func NewPostgresClient(host, user, password, ssl string, port int, logger *Logger) *PostgresClient {
	return &PostgresClient{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		SSLMode:  ssl,
		Logger:   logger,
	}
}

func (psql *PostgresClient) Connect() error {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s", psql.Host, psql.Port, psql.User, psql.Password)
	if psql.DBName != "" {
		psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, psql.DBName)
	}
	if psql.SSLMode != "" {
		psqlInfo = fmt.Sprintf("%s sslmode=%s", psqlInfo, psql.SSLMode)
	}
	//fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	psql.Client = db

	return nil
}

func (psql *PostgresClient) Disconnect() error {
	return psql.Client.Close()
}

func (psql *PostgresClient) Reconnect() error {
	return nil
}

func (psql *PostgresClient) Seed(d DataSet) (interface{}, error) {

	log.Println("Seeding...")

	dsJson := ToJSON(d)
	ds := &DBDataSet{}
	err := json.Unmarshal(dsJson, ds)
	if err != nil {

	}

	ds.Name = sanitize(ds.Name)

	_ = psql.createDB(ds.Name)
	_ = psql.Disconnect()
	psql.DBName = ds.Name
	_ = psql.Connect()
	// createDB(ds.Name)

	tables := ds.Tables
	for n, t := range tables {
		err = psql.createTable(n, &t)
		if err != nil {

		}
	}

	return nil, nil //psql.listDB()
}

func (psql *PostgresClient) Query(s string) (interface{}, error) {

	log.Printf("Query: %s", s)

	result, err := psql.query(s)

	return convertSqlRows(result), err
}

func (psql *PostgresClient) ListDB() (interface{}, error) {

	return psql.listDB()
}

func (psql *PostgresClient) listDB() (interface{}, error) {

	rows, err := psql.query(`SELECT * FROM pg_database;`)

	pdb := PostgresDatabase{}
	pdbs := map[string]PostgresDatabase{}

	for rows.Next() {
		e := rows.Scan(&pdb.Datname, &pdb.Datdba, &pdb.Encoding, &pdb.Datcollate, &pdb.Datctype, &pdb.Datistemplate, &pdb.Datallowconn, &pdb.Datconnlimit, &pdb.Datlastsysoid, &pdb.Datfrozenxid, &pdb.Datminmxid, &pdb.Dattablespace, &pdb.Datacl)
		if e != nil {
			fmt.Println(e)
		}
		//fmt.Println(pdb)
		pdbs[pdb.Datname] = pdb
	}

	// for _, p := range pdbs {
	// 	fmt.Printf("%s\n", p)
	// }

	return pdbs, err
}

func (psql *PostgresClient) newDB(name string) (*sql.Rows, error) {
	return psql.query(fmt.Sprintf("CREATE DATABASE %s", name))
}

func (psql *PostgresClient) query(q string) (*sql.Rows, error) {

	err := psql.Client.Ping()
	if err != nil {
		return nil, err
	}

	rows, err := psql.Client.Query(q)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (psql *PostgresClient) createDB(name string) error {
	query := fmt.Sprintf("CREATE DATABASE %s;", sanitize(name))
	log.Println(query)
	result, err := psql.query(query)
	log.Printf("DB Result: \n%s\n", string(ToJSON(result)))
	log.Printf("DB Error: %s", err)

	return nil
}

func (psql *PostgresClient) dropDB(name string) error {
	query := fmt.Sprintf("DROP DATABASE %s;", sanitize(name))
	log.Println(query)
	result, err := psql.query(query)

	log.Println(result)
	return err
}

func (psql *PostgresClient) dropTable(name string) error {
	query := fmt.Sprintf("DROP TABLE %s;", sanitize(name))
	log.Println(query)
	result, err := psql.query(query)
	log.Printf("Drop table Result: \n%s\n", string(ToJSON(result)))
	log.Printf("Drop table Error: %s", string(ToJSON(err)))

	return err
}

func (psql *PostgresClient) createTable(name string, dt *DataTable) error {

	/*

		CREATE TABLE [IF NOT EXISTS] table_name (
		column1 datatype(length) column_contraint,
		column2 datatype(length) column_contraint,
		column3 datatype(length) column_contraint,
		table_constraints
		);

	*/

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", sanitize(name))
	columns := dt.Columns
	rows := dt.Rows
	delim := dt.Delimiter

	for _, col := range columns {
		query = fmt.Sprintf("%s\n%s %s", query, sanitize(col.Header), col.DataType)
		if len(col.Constraints) > 0 {
			for _, c := range col.Constraints {
				query = fmt.Sprintf("%s %s", query, c)
			}
		}
		query = fmt.Sprintf("%s,", query)
	}

	if len(dt.Constraints) > 0 {
		query = fmt.Sprintf("%s\n", query)
		for _, c := range dt.Constraints {
			query = fmt.Sprintf("%s %s", query, c)
		}
	}

	query = fmt.Sprintf("%s\n);", strings.TrimRight(query, ","))

	log.Println(query)
	log.Println(rows)
	log.Println(delim)

	result, err := psql.query(query)
	log.Println(result)

	return err
}

func sanitize(s string) string {
	lower := strings.ToLower(s)
	noSpace := strings.ReplaceAll(lower, " ", "_")
	return noSpace
}

func convertSqlRows(sr *sql.Rows) *SQLOutput {

	if sr != nil {

		cols, _ := sr.Columns()
		colTypes, _ := sr.ColumnTypes()

		output := &map[string]interface{}{
			"column_02": nil,
			"column_01": nil,
		}
		output2 := []interface{}{}
		output3 := struct {
			column_02 string
			column_01 string
		}{}

		for i, c := range cols {
			//output[c] = nil
			log.Printf("%d: %s", i, c)
		}

		for i, c := range colTypes {
			log.Printf("%d: %s: %s", i, c.Name(), c.ScanType())
		}

		for sr.Next() {
			err := sr.Scan(&output3.column_02, &output3.column_01)
			if err != nil {
				log.Println(err)
			}
			log.Printf("output3: %s", output3)
		}

		log.Printf("output: %s", output)
		log.Printf("output2: %s", output2)

	}

	return nil
}

type SQLOutput struct {
}

// func returnProperties(i interface{}) []interface{} {
// 	for i.
// }
