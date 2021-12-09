package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const stringToSearch = "concurrency"

var sites = []string{
	"https://google.com",
	"https://itc.ua/",
	"https://twitter.com/concurrencyinc",
	"https://twitter.com/",
	"http://localhost:8000",
	"https://github.com/bradtraversy/go_restapi/blob/master/main.go",
	"https://www.youtube.com/",
	"https://postman-echo.com/get",
	"https://en.wikipedia.org/wiki/Concurrency_(computer_science)#:~:text=In%20computer%20science%2C%20concurrency%20is,without%20affecting%20the%20final%20outcome.",
}

type SiteData struct {
	data []byte
	uri  string
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	resultsCh := make(chan SiteData, len(sites))

	findData(ctx, cancel, resultsCh)

	time.Sleep(time.Second)
}

func findData(ctx context.Context, cancel context.CancelFunc, ch chan SiteData) {
	for _, uri := range sites {
		fmt.Println("Sending request to", uri)
		go performRequest(ctx, uri, ch)
	}

	for {
		select {
		case site := <-ch:
			if ok := strings.Contains(string(site.data), stringToSearch); ok {
				fmt.Println(stringToSearch+" string is found in", site.uri)
				cancel()
				return
			}
			fmt.Println("Nothing found in", site.uri)
		}
	}
}

func performRequest(ctx context.Context, uri string, ch chan SiteData) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch <- SiteData{data: bodyBytes, uri: uri}
}
