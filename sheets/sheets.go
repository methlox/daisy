package sheets

import (

	"time"
	"fmt"
	"net/http"
	"os"
	"log"
	"context"
	"database/sql"
	"google.golang.org/api/sheets/v4"
	"golang.org/x/oauth2/google"
	"encoding/base64"
	"google.golang.org/api/option"
)

var db *sql.DB
var err error

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

type Question struct {
	ID				uint		`gorm:"primary key;autoIncrement" json:"id"`
	Form_id			uint 		`json:"form_id"`
	Question_text	*string 	`json:"ques_text"`
	Question_order	*string 	`json:"ques_order"`
	Created_at 		time.Time   `json:"created_at"`
}

type Response struct {
	ID				uint		`gorm:"primary key;autoIncrement" json:"id"`
	Response_id		uint 		`json:"response_id"`
	Question_id		uint 		`json:"ques_id"`
	Answer			*string 	`json:"answer"`
}


func main(config *Config) { 
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Now we are connected to POSTGRESQL DATABASE.")

	http.HandleFunc("/sheets", toSheets)
	http.ListenAndServe(":8080", nil)
}

func toSheets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "We are connected to a browser")
	if r.Method != "GET" {
		http.Error(w, http.StatusText(404), http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT ques_text FROM Question")

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fmt.Fprintf(w, "Great !!! we are connected to a browser\n")
	if r.Method != "GET" {
		http.Error(w, http.StatusText(404), http.StatusMethodNotAllowed)
		return
	}
	col, err := db.Query("SELECT answer FROM Response")

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	defer col.Close()

	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	credBytes, err := base64.StdEncoding.DecodeString(os.Getenv("KEY_JSON_BASE64"))
	if err != nil {
		log.Fatal(err)
		return
	}

	// authenticate and get configuration
	configuration, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatal(err)
		return
	}

	// create client with config and context
	client := configuration.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal(err)
		return
	}

	// write data
	spreadsheetId := "YOUR SPREADSHEET ID"
    writeRange := "A1"
    var vr sheets.ValueRange
    myval := []interface{}{"One", "Two", "Three"}
    vr.Values = append(vr.Values, myval)
    _, err = srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
    if err != nil {
        log.Fatalf("Unable to retrieve data from sheet. %v", err)
    }

}