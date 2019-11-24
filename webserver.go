package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"log"
	"io/ioutil"
	"html/template"
)

// function that renders the correct html template with given data
func renderResponse(w http.ResponseWriter, code int, which string, m map[string]interface{}) {
	// set header values properly
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")	
	// determine the file and render it
	switch which {
	case "upload":
		tmpl := template.Must(template.ParseFiles("static_files/upload.html"))
		tmpl.Execute(w, m)
	case "stats":
		tmpl := template.Must(template.ParseFiles("static_files/stats.html"))
		tmpl.Execute(w, m)
	case "email":
		tmpl := template.Must(template.ParseFiles("static_files/email.html"))
		tmpl.Execute(w, m)
	}
}

// handles the serving of /upload endpoint for GET and POST
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// if request method is GET, we serve the simple static upload form
	if r.Method == "GET" {
		renderResponse(w, 200, "upload", map[string]interface{}{"msg": ""})
    } else if r.Method == "POST" {
		log.Println("File Upload Endpoint Incoming request")

	    // Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	    r.ParseMultipartForm(10 << 20)
	    file, handler, err := r.FormFile("myFile")
	    if err != nil {
	        fmt.Println(err)
	        renderResponse(w, 400, "upload", map[string]interface{}{"msg": "Error receiving file!"})
	    }
	    defer file.Close()

	    // Create a temporary file within our temp-images directory that follows a particular naming pattern
	    fileSave, err := ioutil.TempFile("monitored_folder", "*_"+handler.Filename)
	    if err != nil {
	        fmt.Println(err)
	    }
	    defer fileSave.Close()

	    // read all of the contents of our uploaded file into a byte array
	    fileBytes, err := ioutil.ReadAll(file)
	    if err != nil {
	        fmt.Println(err)
	    }
		
	    // write this byte array to our temporary file
	    fileSave.Write(fileBytes)
		if _, err := fileSave.Write(fileBytes); err != nil {
			log.Println(err)
			renderResponse(w, 400, "upload", map[string]interface{}{"msg": "Something went wrong!"})
		} else {
			// return that we have successfully uploaded our file!
			log.Printf( "Successfully Uploaded File: %s.\n", handler.Filename)
		    renderResponse(w, 201, "upload", map[string]interface{}{"msg": "Upload successful!"})
		}
    } else {
		renderResponse(w, 405, "upload", map[string]interface{}{"msg": "Invalid request type!"})
    }
}

// convenience function that selects the file name with highest compression rate from DB
func selectHighestComp(db *sql.DB) string {
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
	return fileName
}

// convenience function that gets the average compression rate of all files recorded in DB
func selectAvgCompRate(db *sql.DB) float64 {
	rows, err := db.Query("select avg(comp_rate) from data;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

 	var avg float64
	for rows.Next() {
		err = rows.Scan(&avg)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return avg
}

// convenience function that gets the latest 10 logs from DB 
func selectLatestLogs(db *sql.DB) []string {
	rows, err := db.Query("select time, file, comp_rate from data ORDER BY time DESC LIMIT 10;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data []string
 	var avg float64
	var time int 
	var name string
	for rows.Next() {
		err = rows.Scan(&time, &name, &avg)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, fmt.Sprintf("timestamp: %d, compression-rate: %.2f, file: %s", time, avg, name))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// handler function for the /stats endpoing serving basic statistics about the service
func statsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "stats.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	renderResponse(w, 200, "stats", 
			map[string]interface{}{"highestCompRate": selectHighestComp(db),
						 	       "averageCompRate": selectAvgCompRate(db),
						           "lastLogs": selectLatestLogs(db)})
}

// handler function to send email notifications 
// TODO
func emailHandler(w http.ResponseWriter, r *http.Request) {
	renderResponse(w, 200, "email", nil)
}

func main() {
	// basic list of endpoints served by the web interface of the service
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/email", emailHandler)
	http.HandleFunc("/upload", uploadFileHandler)
	http.ListenAndServe(":8080", nil)
}
