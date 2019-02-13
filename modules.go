package main

import (
	"fmt"
	"strings"
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
	resolvedUrl := url.GetUrl()
	fmt.Printf("Fetching url %q\n", url.GetUrl())
	resp, body, errors := gorequest.New().Get(resolvedUrl).End()
	err := squashErrors(errors)
	if err != nil {
		fmt.Printf("Error while fetching %q\n", url.url)
		fmt.Println(err)
		return UrlResponse{}, err
	}
	
	fmt.Printf("Successfully fetched %q. Content Size: %d\n", url.url, len(body))
	return *NewUrlResponse(url.url, resp, body), nil
}

func fetcher(input chan Url, output chan Url) {
	for url := range input {
		urlResponse, error := fetchUrl(url)
		if error == nil {
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
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
  	if err != nil {
    		fmt.Println(err)
		return nil
  	}
	// Process anchor tags. Resolve relative urls. Ignore urls from other host
	a := doc.Find("a")
	fmt.Printf("Url: %q, <a> tags: %d", url, a.Length())
	urlsFromAnchor := a.FilterFunction(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		return exists && IsSameHost(url, href)
	}).Map(func(i int, s *goquery.Selection) string {
		href, _ := s.Attr("href")
		return ResolveUrl(url, href)	
	})
	
	// Handle images. They might come from CDNs, so same host check will fail
	// Intentionally ignoring the other assets from link / script tags
	img := doc.Find("img")
	fmt.Printf("Url: %q, <img> tags: %d", url, img.Length())
	imgUrls := img.FilterFunction(func(i int, s *goquery.Selection) bool {
		_, srcExists := s.Attr("src")
		return srcExists
	}).Map(func(i int, s *goquery.Selection) string {
		src, _ := s.Attr("src")
		return ResolveUrl(url, src)	
	})
	
	outLinks := append(urlsFromAnchor, imgUrls...)
	fmt.Printf("Got %d url(s) from html of %q\n", (len(outLinks)), url)
	return outLinks
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
