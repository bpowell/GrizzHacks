package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
func _main() {
	ArticalUniqeWords(test_string)
}
*/

var convertion_table map[string]string = map[string]string{
	".": "",
	",": "",
	")": "",
	"(": "",
	"-": " ",
	"—": "",
	"“": "",
	"”": "",
}

func ArticalUniqeWords(article string) {
	article = strings.Replace(article, "\n", " ", -1)
	article = strings.ToLower(article)
	for convert_word, convert_replace := range convertion_table {
		article = strings.Replace(article, convert_word, convert_replace, -1)
	}
	words := strings.Split(article, " ")
	fmt.Println(words)
	word_count_map := make(map[string]int)
	for _, word := range words {
		if !strings.Contains(word, " ") && (word != "") {
			word_count_map[word] = word_count_map[word] + 1
		}
	}
	//word_count_map = rankByWordCount(word_count_map)
	rankByWordCount(word_count_map)
	/*
		for key, value := range word_count_map {
			fmt.Println("Key:", key, "Value:", value)
		}
	*/
}

func rankByWordCount(wordFrequencies map[string]int) {
	n := map[int][]string{}
	var a []int
	for k, v := range wordFrequencies {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	for _, k := range a {
		for _, s := range n[k] {
			fmt.Printf("key:%s value:%d\n", s, k)
		}
	}
}
