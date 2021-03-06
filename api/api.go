package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

type ArticleIdAndDate struct {
	Id    int
	Date  string
	Url   string
	Title string
}

type UniqueWords struct {
	Id        int
	Word      string
	Weights   float32
	Count     int
	ArticleId int
}

type NLPWords struct {
	Word   string
	Weight float32
	Ticker string
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
	ticker := strings.ToLower(r.PostFormValue("ticker"))
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
	ticker := strings.ToLower(r.PostFormValue("ticker"))
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

func getDayForStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	ticker := strings.ToLower(r.PostFormValue("ticker"))
	date := r.PostFormValue("date")

	if date == "" || ticker == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	const shortForm = "2006-Jan-02"
	time, _ := time.Parse(shortForm, date)
	start := time.Unix()
	end := time.AddDate(0, 0, 1).Unix()

	stocks := doDatabaseQuery("select id, to_timestamp(timestamp) as timestamp, ticker, close, high, low, open, volume from historic where ticker = $1 and timestamp >= $2 and timestamp <= $3", ticker, start, end)
	if err := json.NewEncoder(w).Encode(stocks); err != nil {
		panic(err)
	}
}

func getAllTickers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rows, err := db.Query("select array_to_json(array_agg(ticker)) as ticker from tickers")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var ticker string
	for rows.Next() {
		if err = rows.Scan(&ticker); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ticker))
}

func getIdsForArticlesForTicker(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	ticker := strings.ToLower(r.PostFormValue("ticker"))

	if ticker == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	var data []ArticleIdAndDate
	rows, err := db.Query("select id, pubdate, url, title from articles where ticker = $1", strings.ToUpper(ticker))
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var d ArticleIdAndDate
		if err = rows.Scan(&d.Id, &d.Date, &d.Url, &d.Title); err != nil {
			panic(err)
		}
		data = append(data, d)
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func getRawArticleById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	id := strings.ToLower(r.PostFormValue("id"))

	if id == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("select raw from articles where id = $1", id)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var raw string
	for rows.Next() {
		if err = rows.Scan(&raw); err != nil {
			panic(err)
		}
	}

	w.Write([]byte(raw))
}

func validateArticleId(id string) error {
	var dbId int
	err := db.QueryRow("select id from articles where id = $1", id).Scan(&dbId)
	switch {
	case err != nil:
		return err
	}

	return nil
}

func updateCountForWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	id := strings.ToLower(r.PostFormValue("id"))
	count, err := strconv.Atoi(r.PostFormValue("count"))

	if id == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	if validateArticleId(id) != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	err = db.QueryRow(`update uniquewords set count = $1 where id = $2 returning id`, count, id).Scan(&id)
	if err != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	w.Write([]byte("OK"))
}

func updateWeightsForWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	id := strings.ToLower(r.PostFormValue("id"))
	weights, err := strconv.ParseFloat(r.PostFormValue("weights"), 32)

	if id == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	if validateArticleId(id) != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	err = db.QueryRow(`update uniquewords set weights = $1 where id = $2 returning id`, weights, id).Scan(&id)
	if err != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	w.Write([]byte("OK"))
}

func addUniqueWordForArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	article_id := strings.ToLower(r.PostFormValue("article_id"))
	word := r.PostFormValue("word")
	weights, err := strconv.ParseFloat(r.PostFormValue("weights"), 32)
	count, err := strconv.Atoi(r.PostFormValue("count"))

	if article_id == "" || word == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	if validateArticleId(article_id) != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	var id int
	err = db.QueryRow(`insert into uniquewords (word, weights, count, article_id) values($1, $2, $3, $4) returning id`, word, weights, count, article_id).Scan(&id)
	if err != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	w.Write([]byte("OK"))
}

func getAllWordsForArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	article_id := strings.ToLower(r.PostFormValue("article_id"))

	if article_id == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	if validateArticleId(article_id) != nil {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("select id, word, weights, count, article_id from uniquewords where article_id = $1", article_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var words []UniqueWords

	for rows.Next() {
		var word UniqueWords
		if err = rows.Scan(&word.Id, &word.Word, &word.Weights, &word.Count, &word.ArticleId); err != nil {
			panic(err)
		}
		words = append(words, word)
	}

	if err := json.NewEncoder(w).Encode(words); err != nil {
		panic(err)
	}
}

func getUniqueWordsInfoByWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	word := strings.ToLower(r.PostFormValue("word"))

	if word == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("select id, weights, count, article_id from uniquewords where word = $1", word)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var things []UniqueWords

	for rows.Next() {
		var t UniqueWords
		if err = rows.Scan(&t.Id, &t.Weights, &t.Count, &t.ArticleId); err != nil {
			panic(err)
		}
		things = append(things, t)
	}

	if err := json.NewEncoder(w).Encode(things); err != nil {
		panic(err)
	}
}

func getNLPWords(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	ticker := strings.ToUpper(r.PostFormValue("ticker"))

	if ticker == "" {
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("select word, uniquewords.weights, articles.ticker from uniquewords, articles where uniquewords.article_id = articles.id and articles.ticker = $1", ticker)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var things []NLPWords

	for rows.Next() {
		var t NLPWords
		if err = rows.Scan(&t.Word, &t.Weight, &t.Ticker); err != nil {
			panic(err)
		}
		things = append(things, t)
	}

	if err := json.NewEncoder(w).Encode(things); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println(config)

	http.HandleFunc("/api/getall", getAllForStock)
	http.HandleFunc("/api/getrange", getRangeForStock)
	http.HandleFunc("/api/getday", getDayForStock)
	http.HandleFunc("/api/gettickers", getAllTickers)
	http.HandleFunc("/api/getarticleids", getIdsForArticlesForTicker)
	http.HandleFunc("/api/getarticle", getRawArticleById)
	http.HandleFunc("/api/updatecount", updateCountForWord)
	http.HandleFunc("/api/updateweights", updateWeightsForWord)
	http.HandleFunc("/api/adduniqueword", addUniqueWordForArticle)
	http.HandleFunc("/api/getwodsforarticle", getAllWordsForArticle)
	http.HandleFunc("/api/getinfoforword", getUniqueWordsInfoByWord)
	http.HandleFunc("/api/getnlpwords", getNLPWords)
	http.ListenAndServe(":8080", nil)
}
