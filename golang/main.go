package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type FeedResponseItem struct {
	Title     string `json:"title"`
	Copyright string `json:"copyright"`
	FullUrl   string `json:"fullUrl"`
	ThumbUrl  string `json:"thumbUrl"`
	ImageUrl  string `json:"imageUrl"`
	Date      string `json:"date"`
	PageUrl   string `json:"pageUrl"`
}

func main() {
	resp, err := http.Get("https://peapix.com/bing/feed")
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(1)
	}

	var items []FeedResponseItem
	err = json.Unmarshal(body, &items)
	if err != nil {
		os.Exit(1)
	}

	var waitGroup sync.WaitGroup
	for i, item := range items {
		waitGroup.Add(1)
		go func(item FeedResponseItem, i int, total int) {
			defer waitGroup.Done()
			download(item, i, total)
		}(item, i, len(items))
	}
	waitGroup.Wait()
}

func download(item FeedResponseItem, index int, total int) {
	resp, err := http.Get(item.ImageUrl)
	if err != nil {
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(1)
	}

	imageUrl, err := url.Parse(item.ImageUrl)
	if err != nil {
		os.Exit(1)
	}

	target := strings.Replace(imageUrl.Path, "/", "", 1)
	err = os.WriteFile(target, body, 0600)
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("%d/%d %s downloaded\n", index, total, item.ImageUrl)
}
