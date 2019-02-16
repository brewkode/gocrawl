package main

import (
	"fmt"
	"strings"
	"regexp"
	"github.com/parnurzeal/gorequest"
	"github.com/hashicorp/go-multierror"
	"github.com/PuerkitoBio/goquery"
)

func shouldCrawl(resolvedUrl string, alreadySeen map[string]bool) bool {
	if _, isPresent := alreadySeen[resolvedUrl]; !isPresent && resolvedUrl != "" {
		// fmt.Printf("Url %q not seen. Forwarding it.\n", resolvedUrl)
		return true
	} else if isPresent {
		// fmt.Printf("Url %q already crawled. Skipping it.\n", resolvedUrl)
		return false
	} else {
		// fmt.Printf("Empty or null url. Ignoring it\n")
		return false
	}
}


var inValidUrlExtensions = regexp.MustCompile(".*(css|js|bmp|gif|jpe?g|svg|png|pnj|mng|tiff?|mid|mp2|mp3|mp4|wav|avi|mov|mpeg|ram|m4v|pdf|wmv|swf|wma|zip|rar|gz|xml|ico|fla|flv|swt|swc)$")
func validUrlExtension(url string) bool {
	return !inValidUrlExtensions.MatchString(url)
}

func urlFilter(input <-chan Url, output chan Url) {
	alreadySeen := make(map[string]bool)
	for url := range input { 
		resolvedUrl := url.GetUrl()
		if shouldCrawl(resolvedUrl, alreadySeen) && validUrlExtension(resolvedUrl) {
			alreadySeen[resolvedUrl] = true
			output <- url
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
	
	fmt.Printf("Successfully fetched %q. StatusCode: %d, Content Size: %d\n", url.url, resp.StatusCode, len(body))
	return *NewUrlResponse(url.url, resp, body), nil
}

func fetcher(input chan Url, output chan Url) {
	for url := range input {
		urlResponse, error := fetchUrl(url)
		if error == nil {
			out := Url{url: urlResponse.url, html: urlResponse.html, redirectedToUrl: urlResponse.redirectedUrl}
			output <- out
		} 
	}
}

func parse(url string, html string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
  	if err != nil {
    		fmt.Println(err)
		return nil
  	}
	// Process anchor tags. Resolve relative urls. Ignore urls from other host
	a := doc.Find("a")
	fmt.Printf("Url: %q, <a> tags: %d\n", url, a.Length())
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
	fmt.Printf("Url: %q, <img> tags: %d\n", url, img.Length())
	imgUrls := img.FilterFunction(func(i int, s *goquery.Selection) bool {
		_, srcExists := s.Attr("src")
		return srcExists
	}).Map(func(i int, s *goquery.Selection) string {
		src, _ := s.Attr("src")
		return ResolveUrl(url, src)	
	})
	
	outLinks := append(urlsFromAnchor, imgUrls...)
	fmt.Printf("Got %d outgoing link(s) from html of %q\n", (len(outLinks)), url)
	return outLinks
}

func linkExtractor(input chan Url, output chan Url, sitemapInput chan Url) {
	for url := range input {
		outLinks := parse(url.GetUrl(), url.html)
		
		// send to site map builder
		go func() {
			sitemapInput <- Url{url: url.GetUrl(), outLinks: outLinks}
		}()
		
		// send back the urls for crawl
		for _, outLink := range outLinks {
			output <- Url{url: outLink}
		}
	}
}

func siteMapBuilder(extractedLinks chan Url, sitemapRequest chan string) {
	var siteMap SiteMap
	for { 
		select {
		case url := <-extractedLinks:
			siteMap.Add(url.url, url.outLinks...)
		case req := <-sitemapRequest:
			fmt.Printf("Sitemap of url(%q) ::\n", req)
			fmt.Printf(siteMap.print(req, 3))
		default:
		}
	}
}

// ----------------- Utilities
func squashErrors(errors []error) error {
	var err *multierror.Error
	for _, e := range errors {
		err = multierror.Append(err, e)
	}
	return err.ErrorOrNil()
}
