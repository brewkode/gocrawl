package main

import (
	"fmt"
	"time"
	"math/rand"
	// "github.com/jinzhu/copier"
	"github.com/parnurzeal/gorequest"
	"github.com/hashicorp/go-multierror"
	"github.com/PuerkitoBio/goquery"
)

func urlFilter(input <-chan Url, output chan Url) {
	alreadySeen := make(map[string]bool)
	for url := range input { 
		resolvedUrl := url.GetUrl()
		if _, isPresent := alreadySeen[resolvedUrl]; !isPresent && resolvedUrl != "" {
			fmt.Printf("Url %q not seen. Forwarding it.\n", resolvedUrl)
			alreadySeen[resolvedUrl] = true
			output <- url
		} else if isPresent {
			fmt.Printf("Url %q already crawled. Skipping it.\n", resolvedUrl)
		} else {
			fmt.Printf("Empty or null url. Ignoring it\n")
		}
	}
}

func fetchUrl(url Url) (UrlResponse, error) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	resolvedUrl := url.GetUrl()
	fmt.Printf("Fetching url %q\n", url.GetUrl())
	resp, body, errors := gorequest.New().Get(resolvedUrl).End()
	err := squashErrors(errors)
	if err != nil {
		fmt.Printf("Error while fetching %q\n", url.url)
		fmt.Println(err)
		return UrlResponse{}, err
	}
	
	fmt.Printf("Successfully fetched %q \n", url.url)
	return *NewUrlResponse(url.url, resp, body), nil
}

func fetcher(input chan Url, output chan Url) {
	for url := range input {
		urlResponse, error := fetchUrl(url)
		if error != nil {
			out := Url{url: urlResponse.url, html: urlResponse.html, redirectedToUrl: urlResponse.redirectedUrl}
			output <- out
		} else {
			// empty body for downstream processing to ignore
			out := Url{url: urlResponse.url}
			output <- out
		}
	}
}

func processHtml(url string, html string) []string {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	doc, err := goquery.NewDocumentFromReader(res.Body)
  	if err != nil {
    		fmt.Println(err)
  	}
	doc.Find("a").FilterFunction(func(i int, s *Selection) bool {
		href, exists := s.Attr("href")
		// TODO: Complete the same host check
	}).Map(func(i int, s *Selection) string {
		// TODO: Resolve url
	})
	
	// Handle images. They might come from CDNs, so same host check will fail
	// Intentionally ignoring the other assets from link / script tags
	doc.Find("img").FilterFunction(func(i int, s *Selection) bool {
		src, srcExists := s.Attr("src")
		href, hrefExists := s.Attr("href")
		return srcExists || hrefExists
	}).Map(func(i int, s *Selection) string {
		// TODO: Resolve url
	})
	fmt.Printf("Got %d url(s) from html of %q\n", 1, url)
	return append(make([]string, 1), url)
}

func linkExtractor(input chan Url, output chan Url) {
	for url := range input {
		for _, outLink := range processHtml(url.GetUrl(), url.html) {
			output <- Url{url: outLink}
		}
	}
}

// ----------------- Utility stuff
func squashErrors(errors []error) error {
	var err *multierror.Error
	for _, e := range errors {
		err = multierror.Append(err, e)
	}
	return err.ErrorOrNil()
}
