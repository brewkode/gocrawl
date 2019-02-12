package main

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/jinzhu/copier"
)


func urlFilter(input <-chan Url, output chan Url) {
	alreadySeen := make(map[string]bool)
	for url := range input { 
		resolvedUrl := url.GetUrl()
		if _, isPresent := alreadySeen[resolvedUrl]; !isPresent {
			fmt.Printf("Url %q not seen. Forwarding it.\n", resolvedUrl)
			alreadySeen[resolvedUrl] = true
			output <- url
		} else {
			fmt.Printf("Url %q already crawled. Skipping it.\n", resolvedUrl)
		}
	}
}

func fetchUrl(url Url) string {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	fmt.Printf("Fetching url %q\n", url.GetUrl())
	return "html page of url: " + url.GetUrl()
}

func fetcher(input chan Url, output chan Url) {
	for url := range input {
		body := fetchUrl(url)
		out := Url{}
		copier.Copy(&out, url)
		out.html = body
		output <- out
	}
}

func processHtml(url string, html string) []string {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	fmt.Printf("Got %d urls from html of %q\n", 1, url)
	return append(make([]string, 1), url)
}

func linkExtractor(input chan Url, output chan Url) {
	for url := range input {
		for _, outLink := range processHtml(url.GetUrl(), url.html) {
			output <- Url{url: outLink}
		}
	}
}
