package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	BUFFER_SIZE := 1024
	// safety buffer value used for all channels - to prevent from deadlock or rogue pages, etc
	// did not find a need for it during test crawls
	SAFETY_BUFFER := 32 
	homePage := "http://brewkode.com"
	urlInput := make(chan Url, BUFFER_SIZE)
	toCrawl := make(chan Url, SAFETY_BUFFER)
	htmlOutput := make(chan Url, SAFETY_BUFFER)
	extractedLinks := make(chan Url, SAFETY_BUFFER)
	siteMapQuery := make(chan string, SAFETY_BUFFER)

	// Seeding
	go func() {
		urlInput <- Url{url: homePage}
	}()

	go urlFilter(urlInput, toCrawl)
	go fetcher(toCrawl, htmlOutput)
	// fans-out extracted links to two different channels
	go linkExtractor(htmlOutput, urlInput, extractedLinks)
	go siteMapBuilder(extractedLinks, siteMapQuery)
	
	fmt.Printf("Please enter the url you want to query its adjacent nodes\n")
	for {
		reader := bufio.NewReader(os.Stdin)
		inpUrlWithNL, _ := reader.ReadString('\n')
		inpUrl := strings.TrimSuffix(inpUrlWithNL, "\n")
		fmt.Printf("Querying for adjacent urls of %q", inpUrl)
		siteMapQuery <- inpUrl
	}
}
