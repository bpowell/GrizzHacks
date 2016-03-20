package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"../classification"
	"../normalize"
)

type ArticleIdAndDate struct {
	Id   int
	Date string
	Url  string
}

var current_file_lenght, current_file int

func main() {

	tickers := []string{"GOOG", "GOOGL", "APPL"}

	for _, ticker_symbol := range tickers {
		article_data := GetArticleId(ticker_symbol)

		current_file_lenght = len(article_data)

		for i, info := range article_data {
			current_file = i
			parsed_article := normalize.GetArticles(strconv.Itoa(info.Id))
			word_distribution := normalize.ArticleUniqeWords(parsed_article)

			fmt.Println(word_distribution)
			var hour_add_on bool = false
			for key, _ := range word_distribution {
				if key == "pm" {
					fmt.Println(key)
					hour_add_on = true
				}
				if key == "am" {
					fmt.Println(key)
					hour_add_on = false
				}
			}
			var article_time string
			for key, _ := range word_distribution {
				if strings.Contains(key, ":") {
					time := strings.Split(key, ":")
					if Hour, err := strconv.Atoi(time[0]); err == nil {
						if Minute, err := strconv.Atoi(time[1]); err == nil {
							if (!hour_add_on && Hour > 7) || (hour_add_on && Hour < 4 && Hour > 0) {
								if hour_add_on {
									article_time = strconv.Itoa(Hour + 12)
								} else {

									article_time = strconv.Itoa(Hour)
								}
								article_time += ":"
								article_time += strconv.Itoa(Minute)
								article_time += ":00"
							}
						}
					}
				}
			}
			new_date := strings.Replace(info.Date, "-", "/", -1)
			date_time := new_date + "-" + article_time
			fmt.Println(date_time)
			weight, _ := classification.ArticleClassifacation(date_time, "0h20m0s")
			fmt.Println(weight)
			AddWordsCountsWeights(info.Id, word_distribution, weight)
		}
	}
}

func AddWordsCountsWeights(id int, words map[string]int, weights float32) {
	var count, total int
	total = len(words)
	count = 0
	for key, value := range words {
		count++
		api_url := "http://104.131.18.185:8080/api/adduniqueword"
		data := url.Values{}
		data.Add("article_id", strconv.Itoa(id))
		data.Add("word", key)
		data.Add("count", strconv.Itoa(value))
		data.Add("weights", strconv.FormatFloat(float64(weights), 'g', -1, 32))

		client := &http.Client{}
		r, _ := http.NewRequest("POST", api_url, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		response, err := client.Do(r)
		if err != nil {
			fmt.Println(err)
			fmt.Println(response.Status)
		}
		fmt.Println(current_file, "/", current_file_lenght, "-", count, "/", total)
	}
}

func GetArticleId(ticker string) []ArticleIdAndDate {
	api_url := "http://104.131.18.185:8080/api/getarticleids"
	data := url.Values{}
	data.Add("ticker", ticker)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", api_url, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	response, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var article_id_and_date []ArticleIdAndDate
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	err = decoder.Decode(&article_id_and_date)
	if err != nil {
		fmt.Println(err)
	}

	return article_id_and_date
	//GetArticles(strconv.Itoa(article_id_and_date[0].Id))
}
