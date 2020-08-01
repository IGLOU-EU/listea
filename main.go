// Copyright 2020 Iglou.eu
// license that can be found in the LICENSE file

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"git.iglou.eu/Laboratory/listea/icon"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

// AutoGenerated struct from standard v1 Gitea API
type AutoGenerated []struct {
	ID      int    `json:"id"`
	HTMLURL string `json:"html_url"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Labels  []struct {
		Name string `json:"name"`
	} `json:"labels"`
	State       string    `json:"state"`
	IsLocked    bool      `json:"is_locked"`
	Comments    int       `json:"comments"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClosedAt    time.Time `json:"closed_at"`
	DueDate     time.Time `json:"due_date"`
	PullRequest struct {
		Merged bool `json:"merged"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// APIResultList inline type def
type APIResultList struct {
	entry []APIResult
}

// APIResult inline type def
type APIResult struct {
	ID         int
	Title      string
	Body       string
	Comments   int
	Repository string
	HTMLURL    string
	LabelsName string
	state      string
	IsLocked   bool
	UpdatedAt  time.Time
	ClosedAt   time.Time
	DueDate    time.Time
	PRMerged   bool
}

// ConfigAPI inline type def
type ConfigAPI struct {
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

// APIRequest inline type def
type APIRequest struct {
	URL  []string
	Hash []string
}

// MenuList for sub instance
type MenuList struct {
	ItemTray []*systray.MenuItem
}

var apiRequest APIRequest
var menuItemList []MenuList

func init() {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), "listea")
	case "darwin":
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "listea")
	case "linux", "freebsd", "netbsd", "openbsd":
		x := os.Getenv("XDG_CONFIG_HOME")

		if x == "" {
			configDir = filepath.Join(os.Getenv("HOME"), ".config", "listea")
		} else {
			configDir = filepath.Join(x, "listea")
		}
	default:
		log.Fatal("Your operating system is not supported")
	}

	configFile := filepath.Join(configDir, "config")

	if !fileExist(configDir) {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !fileExist(configFile) {
		var itr ConfigAPI
		itr.APIURL = "https://gitea.com/api/v1"
		itr.List = make([]List, 1)
		itr.List[0].QueryKey = QueryKey{"", "issues", "open", "", ""}

		json, err := json.MarshalIndent(itr, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(configFile, json, 0750)
		if err != nil {
			log.Fatal(err)
		}
	}

	cf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config ConfigAPI
	err = json.Unmarshal(cf, &config)
	if err != nil {
		log.Fatal(err)
	}

	if config.APIURL == "" {
		log.Fatal("No api Url on config file: ", configFile)
	}

	if config.APIToken == "" {
		log.Fatal("No api Token on config file: ", configFile)
	}

	if len(config.List) < 1 || config.List[0].APIRequest == "" {
		log.Fatal("Empty api request List or no Request configured: ", configFile)
	}

	apiRequest.URL = buildAPIRequest(config)
	apiRequest.Hash = make([]string, len(apiRequest.URL))

	menuItemList = make([]MenuList, len(apiRequest.URL))

	go systray.Run(onReady, nil)
}

func main() {
	for {
		apiResult := make([]APIResultList, len(apiRequest.URL))

		proceedAPIRequest(apiResult[:])
		renderAPISystray(apiResult, menuItemList[:])

		time.Sleep(1 * time.Minute)
	}
}

func fileExist(f string) bool {
	_, err := os.Stat(f)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Listea")
	systray.SetTooltip("Task list viewer with a cup of tea")
}

func renderAPISystray(d []APIResultList, s []MenuList) {
	if len(d) < 1 {
		return
	}

	for i, v := range d {
		if len(v.entry) < 1 {
			continue
		}

		for j := range s[i].ItemTray {
			s[i].ItemTray[j].Hide()
		}

		s[i] = menuItemAPI(v.entry)
	}
}

func menuItemAPI(d []APIResult) MenuList {
	var m MenuList

	s := systray.AddMenuItem(fmt.Sprintf("📋 %s", d[0].Repository), "")
	// For grey style
	s.Disable()
	// To be remove with other entry
	m.ItemTray = append(m.ItemTray, s)

	for _, v := range d {
		t := s.AddSubMenuItem(fmt.Sprintf("    {%d} (%s) %s", v.Comments, v.LabelsName, v.Title), v.Body)

		if v.state == "closed" || v.IsLocked || v.PRMerged {
			t.Disable()
		}

		m.ItemTray = append(m.ItemTray, t)
		go trayIsClicked(m.ItemTray[len(m.ItemTray)-1], v.HTMLURL)
	}

	return m
}

func trayIsClicked(t *systray.MenuItem, url string) {
	for {
		select {
		case <-t.ClickedCh:
			open.Run(url)
		}
	}
}

func proceedAPIRequest(o []APIResultList) {
	for i, v := range apiRequest.URL {
		get, err := http.Get(v)
		if err != nil {
			log.Println("Unable to get on : " + v)
			log.Println(err)
			continue
		}
		defer get.Body.Close()

		if get.StatusCode != 200 {
			log.Println("Can't access to API by:", strings.Split(v, "?")[0])
			continue
		}

		body, err := ioutil.ReadAll(get.Body)
		if err != nil {
			log.Println("Unable read body : " + v)
			log.Println(err)
			continue
		}

		h := sha256.New()
		h.Write([]byte(body))
		hsum := fmt.Sprintf("%x", h.Sum(nil))

		if apiRequest.Hash[i] != hsum {
			o[i] = proceedAPIResult(body)
			apiRequest.Hash[i] = hsum
		}
	}
}

func proceedAPIResult(body []byte) APIResultList {
	var res APIResultList
	var api AutoGenerated

	err := json.Unmarshal(body, &api)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range api {
		var data APIResult

		data.ID = v.ID
		data.Title = v.Title
		data.Body = v.Body
		data.Comments = v.Comments
		data.Repository = v.Repository.FullName
		data.HTMLURL = v.HTMLURL
		data.state = v.State
		data.IsLocked = v.IsLocked
		data.UpdatedAt = v.UpdatedAt
		data.ClosedAt = v.ClosedAt
		data.DueDate = v.DueDate
		data.PRMerged = v.PullRequest.Merged

		if len(v.Labels) > 0 {
			data.LabelsName = v.Labels[0].Name
		}

		res.entry = append(res.entry, data)
	}

	return res
}

func buildAPIRequest(config ConfigAPI) []string {
	var list []string

	for _, v := range config.List {
		list = append(list,
			fmt.Sprintf(
				"%s%s?%stoken=%s",
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

	if query.Q != "" {
		o += fmt.Sprintf("q=%s&", query.Q)
	}

	if query.Type != "" {
		o += fmt.Sprintf("type=%s&", query.Type)
	}

	if query.State != "" {
		o += fmt.Sprintf("state=%s&", query.State)
	}

	if query.Labels != "" {
		o += fmt.Sprintf("labels=%s&", query.Labels)
	}

	if query.Milestones != "" {
		o += fmt.Sprintf("milestones=%s&", query.Milestones)
	}

	return o
}
