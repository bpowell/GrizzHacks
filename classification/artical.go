package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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

func main() {
	ArticleClassifacation("3/18/2016-9:35:00", "6h24m0s")
}

//Format expected MM/DD/YYYY-HH:MM:SS, #h#m#s
func ArticleClassifacation(article_date_time string, interval_time string) (float32, error) {
	time_zone, _ := time.LoadLocation("America/New_York")

	// Start Article time----------------------------
	split_article_date_time := strings.Split(article_date_time, "-")
	split_article_date := strings.Split(split_article_date_time[0], "/")
	split_article_time := strings.Split(split_article_date_time[1], ":")

	article_year, err := strconv.Atoi(split_article_date[2])
	if err != nil {
		return 0, err
	}
	article_month, err := strconv.Atoi(split_article_date[0])
	if err != nil {
		return 0, err
	}
	article_day, err := strconv.Atoi(split_article_date[1])
	if err != nil {
		return 0, err
	}
	article_hour, err := strconv.Atoi(split_article_time[0])
	if err != nil {
		return 0, err
	}
	article_minute, err := strconv.Atoi(split_article_time[1])
	if err != nil {
		return 0, err
	}
	article_second, err := strconv.Atoi(split_article_time[2])
	if err != nil {
		return 0, err
	}

	article_time := time.Date(article_year, time.Month(article_month), article_day, article_hour, article_minute, article_second, 00, time_zone)
	// End Article time----------------------------

	market_open := time.Date(article_year, time.Month(article_month), article_day, 9, 30, 00, 00, time_zone)
	market_close := time.Date(article_year, time.Month(article_month), article_day, 16, 00, 00, 00, time_zone)

	interval, err := time.ParseDuration(interval_time)
	negitive_interval, err := time.ParseDuration("-" + interval_time)
	if err != nil {
		return 0, err
	}

	future_time := article_time.Add(interval)
	past_time := article_time.Add(negitive_interval)

	if future_time.After(market_open) && future_time.Before(market_close) {
		//TODO:Refine this so that i can get a more accurate stock price relitive to the time
		start_amount_stock, end_amount_stock, _ := RetriveStockTick(past_time, future_time)

		starting_close := start_amount_stock.Close
		ending_close := end_amount_stock.Close

		fmt.Println(start_amount_stock, end_amount_stock)
		percent_change := ((ending_close - starting_close) / starting_close) * 100
		fmt.Println(percent_change)
		fmt.Println(start_amount_stock.Open)
		fmt.Println(end_amount_stock.Close)
	} else {
		return 0, errors.New("Your time is outside the bounds of the market")
	}
	return 00, nil
}

func RetriveStockTick(start, end time.Time) (Stock, Stock, error) {
	api_url := "http://104.131.18.185:8080/api/getrange"
	data := url.Values{}
	data.Add("ticker", "ADPNW")
	data.Add("start", strconv.FormatInt(start.Unix(), 10))
	data.Add("end", strconv.FormatInt(end.Unix(), 10))

	fmt.Println(data)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", api_url, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	response, err := client.Do(r)
	if err != nil {
		return Stock{}, Stock{}, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Stock{}, Stock{}, err
	}

	fmt.Println(response.Status)
	var stocks []Stock
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	err = decoder.Decode(&stocks)
	if err != nil {
		return Stock{}, Stock{}, err
	}

	//TODO: Catch the size of the stock array and do checking
	fmt.Println(start, end)
	fmt.Println(start.Unix(), end.Unix())
	fmt.Println(stocks)

	return stocks[0], stocks[len(stocks)-1], nil
}
