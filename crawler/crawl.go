package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/net/html"
)

var config Configuration
var db *sql.DB

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

type Article struct {
	Id            int
	PublishedDate string
	RawArticle    string
	ParsedArticle string
	Ticker        string
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

func main() {
	url := fmt.Sprintf(URL, "GOOGL", "2016-03-18")
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	page := string(body)

	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		panic(err)
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
										fmt.Println(li.FirstChild.Attr)
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
}
