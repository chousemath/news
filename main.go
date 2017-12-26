package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Key represents my API keys
type Key struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Target      string `json:"target"`
}

// ArticleSearch represents a single article search api query
type ArticleSearch struct {
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

// ArticleSearchResult represents a article search api query result
type ArticleSearchResult struct {
	Status    string `json:"status"`
	Copyright string `json:"copyright"`
}

func main() {
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

	queryArticleSearch := ArticleSearch{
		Query:       "trump",
		FilterQuery: "",
		BeginDate:   "",
		EndDate:     "",
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
		fmt.Println(*articleSearchResult)
	} else {
		fmt.Println("Unsuccessful request to Article Search")
	}
}

func makeQueryArticleSearch(key Key, query ArticleSearch) (*http.Response, error) {
	var buffer bytes.Buffer
	buffer.WriteString(key.Target)
	buffer.WriteString("?api-key=")
	buffer.WriteString(key.Value)

	if len(query.Query) > 0 {
		buffer.WriteString("&q=")
		buffer.WriteString(query.Query)
	}
	fmt.Println(buffer.String())
	return http.Get(buffer.String())
}

func getArticles(body []byte) (*ArticleSearchResult, error) {
	// create pointer to ArticleSearchResult
	var a = new(ArticleSearchResult)
	err := json.Unmarshal(body, &a)
	if err != nil {
		fmt.Println("Error parsing articles into struct:", err.Error())
		os.Exit(1)
	}
	return a, err
}
