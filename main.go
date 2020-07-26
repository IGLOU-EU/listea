// Copyright 2020 Iglou.eu
// license that can be found in the LICENSE file

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Config inline type def
type Config struct {
	APIURL   string `json:"api_url"`
	APIToken string `json:"api_token"`
	List     []struct {
		APIRequest string `json:"api_request"`
		QueryKey   struct {
			Q          string `json:"q"`
			Type       string `json:"type"`
			State      string `json:"state"`
			Labels     string `json:"labels"`
			Milestones string `json:"milestones"`
		} `json:"query_key"`
	} `json:"list"`
}

var apiRequest []string

func init() {
	cf, err := ioutil.ReadFile("config")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = json.Unmarshal(cf, &config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)
}

func main() {

}
