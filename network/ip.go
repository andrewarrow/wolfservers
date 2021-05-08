package network

import (
	"fmt"
	"net/http"
	"time"
)

func DoIpGet(ip string) string {
	agent := "agent"

	urlString := fmt.Sprintf("http://%s:8081/hi", ip)
	request, _ := http.NewRequest("GET", urlString, nil)
	request.Header.Set("User-Agent", agent)
	request.Header.Set("Content-Type", "application/json")
	//request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat))
	client := &http.Client{Timeout: time.Second * 500}
	return DoHttpRead("GET", "hi", client, request)
}
