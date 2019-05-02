package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

const (
	dbhost = "localhost"
	dbport = "5432"
	dbuser = "postgres"
	dbpass = "admin"
	dbname = "test"
)

func main() {
	http.HandleFunc("/api/table", createUserHandler)
	http.HandleFunc("/api/user", getTableHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	//...
}

func getTableHandler(w http.ResponseWriter, r *http.Request) {
	//...
}
func initDb() {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}