package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func resetDatabase() error {
	_, err := db.Exec(`DROP TABLE data`)

	return err
}

func setupDatabase() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS data (
		key TEXT PRIMARY KEY,
		val TEXT
	)
	`)

	return err
}

func setVal(key, val string) error {
	_, err := db.Exec(`
	INSERT INTO data (key, val) VALUES ($1, $2)
	ON CONFLICT (key)
		DO UPDATE SET val = $2
	`, key, val)

	return err
}

func getVal(key string) (val string, err error) {
	err = db.QueryRow(`SELECT val FROM data WHERE key = $1`, key).Scan(&val)
	return val, err
}

type req struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Val    string `json:"val"`
}

type resp struct {
	Result string `json:"result,omitempty"`
	Val    string `json:"val,omitempty"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintln(w, "use POST / to get or set keys from the database")
	}

	var req req
	var resp resp
	var err error

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		fmt.Fprintln(w, err)
		return
	}

	if req.Action == "get" {
		resp.Val, err = getVal(req.Key)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	} else if req.Action == "set" {
		if err := setVal(req.Key, req.Val); err != nil {
			fmt.Fprintln(w, err)
			return
		}

		resp.Result = "success"
	} else {
		fmt.Fprintln(w, "unknown action")
	}

	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	if err := encoder.Encode(&resp); err != nil {
		fmt.Fprintln(w, err)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	psqlURL := os.Getenv("DATABASE_URL")
	if psqlURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	var err error

	db, err = sql.Open("postgres", psqlURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	if err := setupDatabase(); err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)

	fmt.Println("Listening on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
