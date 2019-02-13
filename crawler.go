package main

import (
	"bufio"
	"fmt"
	"os"
	//"strconv"
	//"time"
	//"math/rand"
)

func main() {
	BUFFER_SIZE := 1024
	homePage := "http://brewkode.com"
	urlInput := make(chan Url, BUFFER_SIZE)
	toCrawl := make(chan Url)
	htmlOutput := make(chan Url)
	// extractedLinks := make(chan Url)
	
	// Seeding
	go func() {
		urlInput <- Url{url: homePage}
	}()

	go urlFilter(urlInput, toCrawl)
	go fetcher(toCrawl, htmlOutput)
	go linkExtractor(htmlOutput, urlInput)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Printf(text)
}
