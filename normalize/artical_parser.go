package normalize

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var remove_table map[string]string = map[string]string{
	"<script>": "</script>",
	"/*":       "*/",
	"{":        "}",
}

var remove_line []string = []string{
	"<script>",
	"</script>",
	"<\\/script>",
	"//",
	"* ",
}

func GetArticles(id string) string {
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

	body_split := ParseHtml(string(body))
	//body_split := strings.Split(string(body), " ")

	for remove_open, remove_close := range remove_table {
		body_split = RemoveBlocks(body_split, remove_open, remove_close)
	}
	for _, tag := range remove_line {
		body_split = RemoveLines(body_split, tag)
	}

	//ParseHtml(strings.Join(body_split, " "))

	return strings.Join(body_split, " ")

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
	fmt.Println("Removeing Block from tag " + open_tag + "->" + close_tag)
	var start_index, end_index, line_count int

	length := len(body)
	for i := 0; i < length; i++ {
		if strings.Contains(body[i], open_tag) {
			start_index = i
		}
		if strings.Contains(body[i], close_tag) {
			end_index = i + 1
			line_count += end_index - start_index

			body = append(body[:start_index], body[end_index:]...)
			i = start_index
			length -= (end_index - start_index)
		}
	}
	return body
}

func ParseHtml(raw_html string) []string {
	doc, err := html.Parse(strings.NewReader(raw_html))
	if err != nil {
		fmt.Println(err)
	}
	var string_array []string
	var f func(*html.Node, *html.Node)
	f = func(node *html.Node, parent_node_type *html.Node) {
		if parent_node_type != nil {
			switch parent_node_type.DataAtom {

			case atom.P, atom.Ul, atom.Table, atom.A, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.Tr, atom.Td, atom.Th, atom.Span, atom.Strong, atom.Li, atom.Abbr, atom.Div:
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
	return string_array
}
