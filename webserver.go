package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"log"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "stats.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	rows, err := db.Query("SELECT file FROM data ORDER BY comp_rate DESC LIMIT 1;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var fileName string

	for rows.Next() {
		err = rows.Scan(&fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
		
	fmt.Fprintf(w, "Hello, the MAX comp rate was for %s!\n", fileName)
}

func main() {
	http.HandleFunc("/stats", statsHandler)
	http.ListenAndServe(":8080", nil)
}
