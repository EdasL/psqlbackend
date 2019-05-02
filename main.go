package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dbhost = "localhost"
	dbport = 5432
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
	http.HandleFunc("/api/user", createUserHandler)
	http.HandleFunc("/api/table", getTableHandler)
	http.HandleFunc("/api/login", loginHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
func initDb() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpass, dbname)

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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		sqlQuery := `
		SELECT
			username,
			password
		FROM users
		WHERE username=
		`
		err := r.ParseForm()
		sqlQuery += "'" + r.FormValue("username") + "'"
		rows, err := db.Query(sqlQuery)
		fmt.Println(sqlQuery)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		user := userSummary{}

		for rows.Next() {
			err = rows.Scan(
				&user.Username,
				&user.Password,
			)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		}
		fmt.Println(user)
		if user.Password != r.FormValue("password") {
			http.Error(w, "Wrong password", 400)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		sqlStatement := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)`
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		_, err = db.Exec(sqlStatement, r.FormValue("username"), r.FormValue("password"))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.WriteHeader(http.StatusOK)
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
			username,
			password
		FROM users
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
