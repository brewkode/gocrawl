package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"flag"
	"time"
)

func main() {
	homePagePtr := flag.String("home", "http://brewkode.com", "home page of the site that we want to crawl")
	CHAN_BUFFER_SIZE := flag.Int("buffer-size", 100000, "Buffer size of the various channels")
	flag.Parse()

	urlInput := make(chan Url, *CHAN_BUFFER_SIZE)
	toCrawl := make(chan Url, *CHAN_BUFFER_SIZE)
	htmlOutput := make(chan Url, *CHAN_BUFFER_SIZE)
	extractedLinks := make(chan Url, *CHAN_BUFFER_SIZE)
	siteMapQuery := make(chan string, *CHAN_BUFFER_SIZE)

	// Seeding
	go func() {
		urlInput <- Url{url: *homePagePtr}
	}()

	go urlFilter(urlInput, toCrawl)
	go fetcher(toCrawl, htmlOutput)
	// fans-out extracted links to two different channels
	go linkExtractor(htmlOutput, urlInput, extractedLinks)
	go siteMapBuilder(extractedLinks, siteMapQuery)
	
	// identifying DONE
	// taking an approach to see if there are no items in any of the channels 
	// for a prolonged period of time and then deciding DONE based on it.
	emptyQueues := 0
	running := true
	progress(urlInput, toCrawl, htmlOutput)
	for running {
		select {
			case <-time.Tick(30 * time.Second):
				if progress(urlInput, toCrawl, htmlOutput) {
					emptyQueues++
				} else {
					emptyQueues = 0
				}
				
				if emptyQueues > 5 {
					running = false
					break
				}
		}
	}
	fmt.Printf("Crawling completed.\n")
	fmt.Printf("Please enter the url you want to query its adjacent nodes\n")
	for {
		reader := bufio.NewReader(os.Stdin)
		inpUrlWithNL, _ := reader.ReadString('\n')
		inpUrl := strings.TrimSuffix(inpUrlWithNL, "\n")
		siteMapQuery <- inpUrl
	}
}

func progress(urlInput chan Url, toCrawl chan Url, htmlOutput chan Url) bool {
	urlInpSize := len(urlInput)
	toCrawlSize := len(toCrawl)
	htmlOutputSize := len(htmlOutput)
	fmt.Printf("Len of urlInput: %d\n", urlInpSize)
	fmt.Printf("Len of toCrawl: %d\n", toCrawlSize)
	fmt.Printf("Len of htmlOutput: %d\n", htmlOutputSize)
	return urlInpSize == 0 && toCrawlSize == 0 && htmlOutputSize == 0
}
