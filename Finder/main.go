package main

import (
	client "GoFindSecrets/Client"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"gopkg.in/yaml.v2"
)

type GithubInfo struct {
	SessionToken string
	Username     string
}

type configuration struct {
	Repos    []string `yaml:"repos"`
	Keywords []string `yaml:"keywords"`
}

type ParsedData struct {
	Repo    string `json:"repo,omitempty"`
	Url     string `json:"url,omitempty"`
	Keyword string `json:"keyword,omitempty"`
}

func main() {
	var username string
	var sessionToken string
	fmt.Print("Enter Username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter User Session Token: ")
	fmt.Scanln(&sessionToken)
	githubC, err := client.GithubNewClient(username, sessionToken)
	if err != nil {
		log.Fatal(err)
	}
	var config configuration
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Println(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println(err)
	}
	dataToWrite := ""
	for _, repo := range config.Repos {
		for _, keyword := range config.Keywords {
			query := fmt.Sprintf("repo:%s %s", repo, keyword)
			codeListUrlPage := fmt.Sprintf("https://github.com/search?q=%s&type=code", url.QueryEscape(query))
			codeList, err := githubC.R().Get(codeListUrlPage)
			if err != nil {
				log.Print(err)
				continue
			}
			doc, err := htmlquery.Parse(strings.NewReader(string(codeList.Body())))
			if err != nil {
				continue
			}

			node, err := htmlquery.QueryAll(doc, `//div[contains(@class, "search-title")]`)
			if node == nil || err != nil {
				continue
			}
			for _, item := range node {
				url := htmlquery.FindOne(item, "//a")
				urlValue := htmlquery.SelectAttr(url, "href")
				if urlValue == "" {
					continue
				}
				parsedData := ParsedData{
					Repo:    repo,
					Url:     htmlquery.SelectAttr(url, "href"),
					Keyword: keyword,
				}
				jsonData, _ := json.Marshal(parsedData)
				dataToWrite += string(jsonData) + "\n"
			}
			time.Sleep(1 * time.Second)
		}
		time.Sleep(10 * time.Second)
	}
	if dataToWrite != "" {
		currentTime := time.Now()
		fileName := fmt.Sprintf("%v.json", currentTime.Format("2006-01-02-15-04-05"))
		os.WriteFile(fileName, []byte(dataToWrite), 0644)
	}
}
