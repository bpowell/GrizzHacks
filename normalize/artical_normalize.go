package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	ArticalUniqeWords(test_string)
}

//var blacklist []string = []string{".", ")", "(", ",", "-"}

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

var test_string string = `When Alphabet (GOOGL) released Q4 earnings in February, the tech giant revealed that last year it spent $3.6 billion on so-called moonshot projects ranging from self-driving cars to extending human lifespans.

While the potential payoff for these initiatives is years away, investors may profit sooner if the stock breaks out of its current base.

Taking Stock Of Alphabet’s Initiatives

Alphabet generated the bulk of its $75 billion in 2015 revenue from advertising, with ads accounting for 90% of Q4 sales. But as its decision to put robot maker Boston Dynamics up for sale due to a lack of commercial prospects indicates, the sultan of search is serious about monetizing its initiatives. Here’s a quick overview of selected projects.

Self-Driving Cars: Alphabet wants Congress to grant it expedited permission to bring self-driving cars to market — vehicles with no pedals or steering wheels.  As the San Jose Mercury News reported, the move comes as Alphabet is dropping “increasingly strong hints” that autonomous vehicles may be ready to roll sooner than expected.

Cloud Services: The cloud computing market is expected to reach $27.4 billion this year. Google’s Cloud Platform trails Amazon (AMZN) unit Amazon Web Services, which has approximately 37% market share, and Azure from Microsoft (MSFT).  But it may get a boost from Apple (AAPL), which signed a deal worth between $400 million and $600 million to use Google’s Cloud Platform for its iCloud service.

Robotic-Assisted Surgery: Last year, Alphabet’s Verily unit teamed up with Johnson & Johnson (JNJ) to launch Verb Surgical, a joint venture to develop advanced surgical robotics to compete with Intuitive Surgical (ISRG).`

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
