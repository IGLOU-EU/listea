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
	List     []List `json:"list"`
}

// QueryKey inline type def
type QueryKey struct {
	Q          string `json:"q"`
	Type       string `json:"type"`
	State      string `json:"state"`
	Labels     string `json:"labels"`
	Milestones string `json:"milestones"`
}

// List inline type def
type List struct {
	APIRequest string   `json:"api_request"`
	QueryKey   QueryKey `json:"query_key"`
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

	apiRequest = buildAPIRequest(config)
	fmt.Println(apiRequest)
}

func main() {

}

func buildAPIRequest(config Config) []string {
	var list []string

	for _, v := range config.List {
		list = append(list,
			fmt.Sprintf(
				"%s%s?%stoken=%%20%%20%s",
				config.APIURL,
				v.APIRequest,
				buildAPIQuery(v.QueryKey),
				config.APIToken,
			),
		)
	}

	return list
}

func buildAPIQuery(query QueryKey) string {
	o := ""

	if !empty(query.Q) {
		o += fmt.Sprintf("q=%s&", query.Q)
	}

	if !empty(query.Type) {
		o += fmt.Sprintf("type=%s&", query.Type)
	}

	if !empty(query.State) {
		o += fmt.Sprintf("state=%s&", query.State)
	}

	if !empty(query.Labels) {
		o += fmt.Sprintf("labels=%s&", query.Labels)
	}

	if !empty(query.Milestones) {
		o += fmt.Sprintf("milestones=%s&", query.Milestones)
	}

	return o
}

func empty(v string) bool {
	if v != "" {
		return true
	}

	return false
}
