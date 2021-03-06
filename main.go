// namespace declaration
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/chousemath/news/utilities"
)

// Key represents my API keys
type Key struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Target      string `json:"target"`
}

// Query represents a single article search api query
type Query struct {
	Query       string
	FilterQuery string
	BeginDate   string
	EndDate     string
	Sort        string
	Fields      string
	Highlight   bool
	Page        int
	FacetField  string
	FacetFilter bool
}

// SearchResult represents a article search api query result
type SearchResult struct {
	Status    string                `json:"status"`
	Copyright string                `json:"copyright"`
	Response  ArticleSearchResponse `json:"response"`
}

// Meta is the meta data tag for each Article Search response
type Meta struct {
	Hits   int `json:"hits"`
	Offset int `json:"offset"`
	Time   int `json:"time"`
}

// Media represents a single multimedia item for a doc
type Media struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// Headline represents the headline for the result doc
type Headline struct {
	Main          string `json:"main"`
	Kicker        string `json:"kicker"`
	ContentKicker string `json:"content_kicker"`
	PrintHeadline string `json:"print_headline"`
}

// Keyword represents the keywords of a particular NYT article
type Keyword struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Doc is a single document from the Article Search response
type Doc struct {
	URL             string    `json:"web_url"`
	Snippet         string    `json:"snippet"`
	PrintPage       string    `json:"print_page"`
	Source          string    `json:"source"`
	PublicationDate string    `json:"pub_date"`
	DocumentType    string    `json:"document_type"`
	SectionName     string    `json:"section_name"`
	Multimedia      []Media   `json:"multimedia"`
	Headline        Headline  `json:"headline"`
	Keywords        []Keyword `json:"keywords"`
}

// ArticleSearchResponse is the main data body of the Article Search response
type ArticleSearchResponse struct {
	Meta Meta  `json:"meta"`
	Docs []Doc `json:"docs"`
}

func main() {
	currentTime := time.Now().Local().Format("20060102")
	fmt.Println("Current Time")
	utilities.ColorPrintln(currentTime, "white")
	queryPtr := flag.String("q", "", "Main query value to pass to the NYT API.")
	filterQueryPtr := flag.String("fq", "", "Filter query value (Lucene syntax)")
	beginDatePtr := flag.String("bd", currentTime, "Begin date for news articles")
	flag.Parse()

	if len(*beginDatePtr) > 0 {
		if *beginDatePtr == "yesterday" {
			*beginDatePtr = time.Now().Local().AddDate(0, 0, -1).Format("20060102")
			fmt.Println("Yesterday's date?")
			fmt.Println(*beginDatePtr)
		}
	}

	var articleSearch Key

	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("There was an error reading the config file:", err.Error())
		os.Exit(1)
	}

	var keys []Key

	json.Unmarshal(raw, &keys)

	for _, key := range keys {
		if key.Title == "articleSearch" {
			articleSearch.Target = key.Target
			articleSearch.Value = key.Value
		}
	}
	fmt.Println("keyArticleSearch:", articleSearch.Target)
	fmt.Println("keyArticleSearch:", articleSearch.Value)

	queryArticleSearch := Query{
		Query:       *queryPtr,
		FilterQuery: *filterQueryPtr,
		BeginDate:   *beginDatePtr,
		EndDate:     currentTime,
		Sort:        "",
		Fields:      "",
		Highlight:   false,
		Page:        -1,
		FacetField:  "",
		FacetFilter: false,
	}

	resp, err := makeQueryArticleSearch(articleSearch, queryArticleSearch)

	if err != nil {
		fmt.Println("There was an error making GET request to Article Search:", err.Error())
		os.Exit(1)
	}

	if resp.StatusCode == 200 {
		fmt.Println(resp)
		fmt.Println("-----")
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("There was an error reading the response body for Article Search:", err.Error())
			os.Exit(1)
		}

		articleSearchResult, err := getArticles(body)
		for i, doc := range articleSearchResult.Response.Docs {
			printDoc(i, doc)
		}
	} else {
		fmt.Println("Unsuccessful request to Article Search")
	}
}

func printDoc(index int, doc Doc) {
	utilities.ColorPrintln(fmt.Sprintf("News Article %v", index), "cyan")
	utilities.ColorPrintln(doc.URL, "magenta")
	utilities.ColorPrintln(doc.Snippet, "yellow")
	utilities.ColorPrintln("-------------\n-------------", "red")
}

func makeQueryArticleSearch(key Key, query Query) (*http.Response, error) {
	var buffer bytes.Buffer
	buffer.WriteString(key.Target)
	buffer.WriteString("?api-key=")
	buffer.WriteString(key.Value)

	buffer.WriteString(utilities.AddToQuery("q", query.Query))
	buffer.WriteString(utilities.AddToQuery("fq", query.FilterQuery))

	// begin_date will go into the query string no matter what
	buffer.WriteString("&begin_date=")
	buffer.WriteString(query.BeginDate)
	fmt.Println(buffer.String())
	return http.Get(buffer.String())
}

func getArticles(body []byte) (*SearchResult, error) {
	// create pointer to SearchResult
	var a = new(SearchResult)
	err := json.Unmarshal(body, &a)
	if err != nil {
		fmt.Println("Error parsing articles into struct:", err.Error())
		os.Exit(1)
	}
	return a, err
}
