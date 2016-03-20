package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type ArticleIdAndDate struct {
	Id   int
	Date string
}

var (
	spacingRe = regexp.MustCompile(`[ \r\n\t]+`)
	newlineRe = regexp.MustCompile(`\n\n+`)
)

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

	//fmt.Println(article_id_and_date)
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

	ParseHtml(string(body))
}

func ParseHtml(raw_html string) {
	doc, err := html.Parse(strings.NewReader(raw_html))
	if err != nil {
		fmt.Println(err)
	}
	var string_array []string
	var f func(*html.Node, *html.Node)
	f = func(node *html.Node, parent_node_type *html.Node) {
		fmt.Println(node, parent_node_type)
		if parent_node_type != nil {
			switch parent_node_type.DataAtom {

			case atom.P, atom.Ul, atom.Table, atom.A, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.Tr, atom.Td, atom.Th, atom.Span, atom.Strong, atom.Li:
				switch node.Type {
				case html.TextNode:
					string_array = append(string_array, node.Data)
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c, node)
		}
	}
	f(doc, nil)
	fmt.Println(string_array)
	ArticleUniqeWords(strings.Join(string_array, " "))
	ArticleUniqeWords(strings.Join(string_array, " "))
}
