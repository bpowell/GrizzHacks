package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/net/html"
)

var config Configuration
var db *sql.DB

const shortForm = "2006-Jan-02"

const URL = "http://finance.yahoo.com/q/h?s=%s&t=%s"

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

type Data struct {
	Ticker string
	Date   string
	Urls   []string
	Url    string
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

func workerpoolGetUrlsToGrab(i int, jobs <-chan Data, results chan<- Data) {
	for data := range jobs {
		fmt.Printf("%d working on %s\n", i, data)
		urls, err := getLinks(data)
		if err != nil {
			results <- data
			continue
		}

		data.Urls = urls
		results <- data
	}
}

func getAllTickers() []string {
	rows, err := db.Query("select upper(ticker) from wtickers")
	//rows, err := db.Query("select upper(ticker) from tickers")
	if err != nil {
		panic(err)
	}

	var tickers []string

	for rows.Next() {
		var ticker string
		if err = rows.Scan(&ticker); err != nil {
			panic(err)
		}
		tickers = append(tickers, ticker)
	}

	return tickers
}

func getDates() []string {
	var dates []string

	start, _ := time.Parse(shortForm, "2016-Mar-10")
	end, _ := time.Parse(shortForm, "2016-Mar-18")

	for current := start; !current.Equal(end); current = current.AddDate(0, 0, 1) {
		if current.Weekday() == 0 || current.Weekday() == 6 {
			continue
		}

		date := fmt.Sprintf("%d-%d-%d", current.Year(), current.Month(), current.Day())
		dates = append(dates, date)
	}

	return dates
}

func getLinks(data Data) ([]string, error) {
	var urls []string

	url := fmt.Sprintf(URL, data.Ticker, data.Date)

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return urls, err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Cannot get %s\n", url)
		return urls, errors.New("Bad")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return urls, err
	}

	page := string(body)

	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		return urls, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			tbody := n.FirstChild
			for tr := tbody.FirstChild; tr != nil; tr = tr.NextSibling {
				td := tr.FirstChild
				if td.Type == html.ElementNode && td.Data == "td" {
					div := td.FirstChild
					if div == nil {
						break
					}

					if div.Type == html.ElementNode && div.Data == "div" {
						for c := div.FirstChild; c != nil; c = c.NextSibling {
							if c.Type == html.ElementNode && c.Data == "ul" {
								for li := c.FirstChild; li != nil; li = li.NextSibling {
									if li.FirstChild.Type == html.ElementNode && li.FirstChild.Data == "a" {
										for _, a := range li.FirstChild.Attr {
											urls = append(urls, a.Val)
											title := li.FirstChild.FirstChild.Data
											insertTitle(a.Val, title)
										}
									}
								}
							}
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return urls, nil
}

func insertTitle(url, title string) {
	var id int
	err := db.QueryRow(`update articles set title = $1 where url = $2 returning id`, title, url).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func main() {
	jobs := make(chan Data, 100)
	results := make(chan Data, 100)

	for w := 0; w < 4; w++ {
		go workerpoolGetUrlsToGrab(w, jobs, results)
	}

	var datas []Data

	for _, ticker := range getAllTickers() {
		for _, date := range getDates() {
			var data Data
			jobs <- Data{ticker, date, nil, ""}
			data = <-results
			datas = append(datas, data)
		}
	}
}
