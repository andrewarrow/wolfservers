package network

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func BaseUrl() string {
	url := os.Getenv("WOLF_DO_HOST")
	if url == "" {
		return "https://api.digitalocean.com/"
	}
	return url
}

func DoGet(pat, route string) string {
	agent := "agent"

	urlString := fmt.Sprintf("%s%s", BaseUrl(), route)
	request, _ := http.NewRequest("GET", urlString, nil)
	request.Header.Set("User-Agent", agent)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat))
	client := &http.Client{Timeout: time.Second * 500}
	return DoHttpRead("GET", route, client, request)
}

func DoHttpRead(verb, route string, client *http.Client, request *http.Request) string {
	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("\n\nERROR: %d %s\n\n", resp.StatusCode, err.Error())
			os.Exit(1)
			return ""
		}
		if resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 204 {
			return string(body)
		} else {
			fmt.Printf("\n\nERROR: %d %s\n\n", resp.StatusCode, string(body))
			os.Exit(1)
			return ""
		}
	}
	fmt.Printf("\n\nERROR: %s\n\n", err.Error())
	os.Exit(1)
	return ""
}

func DoPost(pat, route string, payload []byte) string {
	body := bytes.NewBuffer(payload)
	urlString := fmt.Sprintf("%s%s", BaseUrl(), route)
	request, _ := http.NewRequest("POST", urlString, body)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat))
	client := &http.Client{Timeout: time.Second * 50}

	return DoHttpRead("POST", route, client, request)
}
func DoDelete(route string) string {
	body := bytes.NewBuffer([]byte{})
	urlString := fmt.Sprintf("%s%s", BaseUrl(), route)
	request, _ := http.NewRequest("DELETE", urlString, body)
	request.Header.Set("Content-Type", "application/json")
	//request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat))
	client := &http.Client{Timeout: time.Second * 50}

	return DoHttpRead("DELETE", route, client, request)
}
