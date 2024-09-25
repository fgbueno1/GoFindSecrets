package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	ApiKey   string   `yaml:"api-key"`
	Repos    []string `yaml:"repos"`
	Orgs     []string `yaml:"orgs"`
	Keywords []string `yaml:"keywords"`
}

func main() {
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
		if repo == "" {
			continue
		}
		for _, keyword := range config.Keywords {
			query := url.QueryEscape(fmt.Sprintf("%s repo:%s", keyword, repo))
			results := GitSearch(config.ApiKey, query)
			for _, result := range results.Items {
				parsedData := ParsedData{
					Org:     result.Repository.Owner.Login,
					Repo:    result.Repository.FullName,
					File:    result.Name,
					Url:     result.HTMLURL,
					Keyword: keyword,
				}
				jsonData, _ := json.Marshal(parsedData)
				dataToWrite += string(jsonData) + "\n"
			}
			time.Sleep(1 * time.Second)
		}
		time.Sleep(10 * time.Second)
	}
	for _, org := range config.Orgs {
		if org == "" {
			continue
		}
		for _, keyword := range config.Keywords {
			query := url.QueryEscape(fmt.Sprintf("%s org:%s", keyword, org))
			results := GitSearch(config.ApiKey, query)
			for _, result := range results.Items {
				parsedData := ParsedData{
					Org:     result.Repository.Owner.Login,
					Repo:    result.Repository.FullName,
					File:    result.Name,
					Url:     result.HTMLURL,
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

func GitSearch(apiKey string, query string) GitSearchResult {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s", query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var jsonData GitSearchResult
	err = json.Unmarshal(bodyText, &jsonData)
	if err != nil {
		fmt.Print("Error Unmarshal")
		log.Fatal(err)
	}
	return jsonData
}
