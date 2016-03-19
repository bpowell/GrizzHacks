package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var config Configuration
var db *sql.DB

type databaseInfo struct {
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
}

type Configuration struct {
	Db databaseInfo
}

type Stock struct {
	Id        int
	Timestamp string
	Ticker    string
	Close     float32
	High      float32
	Low       float32
	Open      float32
	Volume    float32
}

type StockRequest struct {
	Ticker    string
	StartTime string
	EndTime   string
}

func init() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Db.Username, config.Db.Password, config.Db.Host, config.Db.Dbname)
	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}
}

func doDatabaseQuery(sql string, args ...interface{}) []Stock {
	rows, err := db.Query(sql, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var stocks []Stock

	for rows.Next() {
		var stock Stock
		if err = rows.Scan(&stock.Id, &stock.Timestamp, &stock.Ticker, &stock.Close, &stock.High, &stock.Low, &stock.Open, &stock.Volume); err != nil {
			panic(err)
		}
		stocks = append(stocks, stock)
	}

	return stocks
}

func getAllForStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	ticker := r.PostFormValue("ticker")
	if ticker == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	stocks := doDatabaseQuery("select id, to_timestamp(timestamp) as timestamp, ticker, close, high, low, open, volume from historic where ticker = $1", ticker)

	if err := json.NewEncoder(w).Encode(stocks); err != nil {
		panic(err)
	}
}

func getRangeForStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	ticker := r.PostFormValue("ticker")
	start := r.PostFormValue("start")
	end := r.PostFormValue("end")

	if start == "" || ticker == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	var stocks []Stock
	if end == "" {
		stocks = doDatabaseQuery("select id, to_timestamp(timestamp) as timestamp, ticker, close, high, low, open, volume from historic where ticker = $1 and timestamp >= $2", ticker, start)
	} else {
		stocks = doDatabaseQuery("select id, to_timestamp(timestamp) as timestamp, ticker, close, high, low, open, volume from historic where ticker = $1 and timestamp >= $2 and timestamp <= $3", ticker, start, end)
	}

	if err := json.NewEncoder(w).Encode(stocks); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println(config)

	http.HandleFunc("/api/getall", getAllForStock)
	http.HandleFunc("/api/getrange", getRangeForStock)
	http.ListenAndServe(":8080", nil)
}
