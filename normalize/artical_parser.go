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

	"golang.org/x/net/html"
)

type ArticleIdAndDate struct {
	Id   int
	Date string
}

var remove_table map[string]string = map[string]string{
	"<script>": "</script>",
	"/*":       "*/",
	"{":        "}",
}
var remove_line []string = []string{
	"<script>",
	"</script>",
	"<\\/script>",
	"function",
}

func main() {
	api_url := "http://104.131.18.185:8080/api/getarticleids"
	data := url.Values{}
	data.Add("ticker", "GOOGL")

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
		//return ArticleIdAndDate{}, ArticleIdAndDate{}, err
	}

	fmt.Println(article_id_and_date)
	GetArticles(strconv.Itoa(article_id_and_date[0].Id))
}

func GetArticles(id string) {
	api_url := "http://104.131.18.185:8080/api/getarticle"
	data := url.Values{}
	data.Add("id", id)

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
	body_split := strings.Split(string(body), " ")
	fmt.Println(len(body_split))

	for remove_open, remove_close := range remove_table {
		body_split = RemoveBlocks(body_split, remove_open, remove_close)
	}
	for _, tag := range remove_line {
		body_split = RemoveLines(body_split, tag)
	}
	body_split = RemoveBlocks(body_split, "<script>", "Dink")
	fmt.Println(len(body_split))
	fmt.Println(body_split)
}

func RemoveLines(body []string, tag string) []string {
	fmt.Println("Removeing Lines with tag " + tag)
	length := len(body)
	for i := 0; i < length; i++ {
		if strings.Contains(body[i], tag) {
			body = append(body[:i], body[i+1:]...)
			i--
			length--
		}
	}
	return body
}

func RemoveBlocks(body []string, open_tag, close_tag string) []string {
	fmt.Println("Removeing Lines from tag " + open_tag + "->" + close_tag)
	var start_index, end_index, line_count int

	length := len(body)
	for i := 0; i < length; i++ {
		if strings.Contains(body[i], open_tag) {
			fmt.Println(body[i])
			start_index = i
		}
		if strings.Contains(body[i], close_tag) {
			fmt.Println(body[i])
			end_index = i + 1
			line_count += end_index - start_index
			fmt.Println(start_index, end_index, line_count)

			body = append(body[:start_index], body[end_index:]...)
			i = start_index
			length -= (end_index - start_index)
			fmt.Println(i)
			fmt.Println(len(body))
			fmt.Println(length)
		}
		//fmt.Println(i, length)
	}
	return body
}

func ParseHtml(raw_html string) {
	doc, err := html.Parse(strings.NewReader(raw_html))
	if err != nil {
		fmt.Println(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// Do something with n...
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}
