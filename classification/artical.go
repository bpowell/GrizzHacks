package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	ArticleClassifacation("3/19/2016-15:50:00", "0h20m0s")
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
	if err != nil {
		return 0, err
	}

	future_time := article_time.Add(interval)

	if future_time.After(market_open) && future_time.Before(market_close) {
		fmt.Println(future_time)
		//TODO: request json for date and time+interal
		fmt.Println("During market hour")
	} else {
		return 0, errors.New("Your time is outside the bounds of the market")
	}
	return 00, nil
}
