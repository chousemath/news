package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Key is a struct to represent my API keys
type Key struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

func main() {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("There was an error:", err.Error())
		os.Exit(1)
	}
	var keys []Key
	json.Unmarshal(raw, &keys)
	fmt.Println(keys[0].Value)
}
