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
	"time"

	"../classification"

	"../normalize"
)

type UniqueWords struct {
	Id        int
	Word      string
	Weights   float32
	Count     int
	ArticleId int
}

func main() {

	var correct, incorrect, total int
	tickers := []string{"GOOG", "GOOGL", "APPL"}

	for _, ticker_symbol := range tickers {
		article_data := GetArticleId(ticker_symbol)

		for _, info := range article_data {
			parsed_article := normalize.GetArticles(strconv.Itoa(info.Id))
			word_distribution := normalize.ArticleUniqeWords(parsed_article)

			var hour_add_on bool = false
			for key, _ := range word_distribution {
				if key == "pm" {
					hour_add_on = true
				}
				if key == "am" {
					hour_add_on = false
				}
			}
			var article_time_parse string
			for key, _ := range word_distribution {
				if strings.Contains(key, ":") {
					time := strings.Split(key, ":")
					if Hour, err := strconv.Atoi(time[0]); err == nil {
						if Minute, err := strconv.Atoi(time[1]); err == nil {
							if (!hour_add_on && Hour > 7) || (hour_add_on && Hour < 4 && Hour > 0) {
								if hour_add_on {
									article_time_parse = strconv.Itoa(Hour + 12)
								} else {

									article_time_parse = strconv.Itoa(Hour)
								}
								article_time_parse += ":"
								article_time_parse += strconv.Itoa(Minute)
								article_time_parse += ":00"
							}
						}
					}
				}
			}
			new_date := strings.Replace(info.Date, "-", "/", -1)
			article_date_time := new_date + "-" + article_time_parse
			interval_time := "0h20m0s"

			time_zone, _ := time.LoadLocation("America/New_York")

			// Start Article time----------------------------
			split_article_date_time := strings.Split(article_date_time, "-")
			split_article_date := strings.Split(split_article_date_time[0], "/")
			split_article_time := strings.Split(split_article_date_time[1], ":")
			if len(split_article_time) < 3 {
				continue
			}
			article_year, _ := strconv.Atoi(split_article_date[0])
			article_month, _ := strconv.Atoi(split_article_date[1])
			article_day, _ := strconv.Atoi(split_article_date[2])
			article_hour, _ := strconv.Atoi(split_article_time[0])
			article_minute, _ := strconv.Atoi(split_article_time[1])
			article_second, _ := strconv.Atoi(split_article_time[2])

			article_time := time.Date(article_year, time.Month(article_month), article_day, article_hour, article_minute, article_second, 00, time_zone)
			// End Article time----------------------------

			market_open := time.Date(article_year, time.Month(article_month), article_day, 9, 30, 00, 00, time_zone)
			market_close := time.Date(article_year, time.Month(article_month), article_day, 16, 00, 00, 00, time_zone)

			interval, _ := time.ParseDuration(interval_time)
			negitive_interval, _ := time.ParseDuration("-" + interval_time)

			future_time := article_time.Add(interval)
			past_time := article_time.Add(negitive_interval)

			if future_time.After(market_open) && future_time.Before(market_close) {
				start_amount_stock, end_amount_stock, _ := classification.RetriveStockTick(ticker_symbol, past_time, future_time)

				starting_close := start_amount_stock.Close
				ending_close := end_amount_stock.Close

				percent_change := ((ending_close - starting_close) / starting_close) * 100
				predict := GetPrediciton(word_distribution)
				if (predict > 0 && percent_change > 0) || predict < 0 && percent_change < 0 {
					correct++
					total++
				} else {
					incorrect++
					total++
				}
				accuracy := ((correct - incorrect) / (total)) * 100
				fmt.Println(accuracy)
			} else {

				fmt.Println("Your time is outside the bounds of the market")
			}
			fmt.Println(len(word_distribution))
		}
	}
	fmt.Println(accuracy)
}

func GetPrediciton(words map[string]int) float32 {
	var sum float32
	var total int
	for word, _ := range words {
		api_url := "http://104.131.18.185:8080/api/getinfoforword"
		data := url.Values{}
		data.Add("word", word)

		client := &http.Client{}
		r, _ := http.NewRequest("POST", api_url, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		response, _ := client.Do(r)
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		fmt.Println(response.Status)
		fmt.Println(string(body))
		fmt.Println(word)
		var ws []UniqueWords
		decoder := json.NewDecoder(strings.NewReader(string(body)))
		err := decoder.Decode(&ws)
		if err != nil {
			return 0
		}
		for _, w := range ws {
			sum += (float32(w.Count) * w.Weights)
			total += w.Count
		}
		fmt.Println("SumTotak:", sum, total)
	}
	return (sum / float32(total))
}
