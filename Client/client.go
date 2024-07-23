package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type GithubClient struct {
	*resty.Client
}

// GithubNewClient Creates a new client to interact with github search tool.
func GithubNewClient(username string, sessionToken string) (GithubClient, error) {
	client := resty.New()
	client.SetRetryCount(3)
	client.SetRetryWaitTime(5 * time.Second)

	usrAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"
	userSession := fmt.Sprintf("user_session=%s; ", sessionToken)
	constValues := `color_mode=%7B%22color_mode%22%3A%22auto%22%2C%22light_theme%22%3A%7B%22name%22%3A%22light%22%2C%22color_mode%22%3A%22light%22%7D%2C%22dark_theme%22%3A%7B%22name%22%3A%22dark%22%2C%22color_mode%22%3A%22dark%22%7D%7D; logged_in=yes; `
	user := fmt.Sprintf("dotcom_user=%s;", username)
	cookie := userSession + constValues + user

	client.SetHeader("Accept", "text/html")
	client.SetHeader("Accept-Encoding", "gzip")
	client.SetHeader("Accept-Language", "en-US,en;q=0.9")
	client.SetHeader("Cookie", cookie)
	client.SetHeader("Host", "github.com")
	client.SetHeader("sec-ch-ua", `"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`)
	client.SetHeader("sec-ch-ua-mobile", "?0")
	client.SetHeader("Sec-Fetch-Dest", "empty")
	client.SetHeader("Sec-Fetch-Mode", "cors")
	client.SetHeader("Sec-Fetch-Site", "same-origin")
	client.SetHeader("User-Agent", usrAgent)
	client.SetHeader("X-Requested-With", "XMLHttpRequest")

	return GithubClient{client}, nil
}
