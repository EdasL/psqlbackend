package main

import (
	"database/sql"
	"encoding/json"
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

type newUser struct {
	Username string
	Password string
}
type userSummary struct {
	Username string
	Password string
}

type users struct {
	Repositories []userSummary
}

func main() {
	initDb()
	defer db.Close()
	http.HandleFunc("/api/table", createUserHandler)
	http.HandleFunc("/api/user", getTableHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
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

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		sqlStatement := `
		INSERT INTO Users (Username, Password)
		VALUES ($1, $2)`
		r.ParseForm()

		_, err := db.Exec(sqlStatement, r.FormValue("username"), r.FormValue("password"))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}

func getTableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		users := users{}

		err := queryTable(&users)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		out, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprintf(w, string(out))
	}
}

func queryTable(users *users) error {
	rows, err := db.Query(`
		SELECT
			Username,
			Password
		FROM Users
		`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		user := userSummary{}
		err = rows.Scan(
			&user.Username,
			&user.Password,
		)
		if err != nil {
			return err
		}
		users.Repositories = append(users.Repositories, user)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	fmt.Println(users)
	return nil
}
