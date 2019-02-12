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
	homePage := "http://www.example.com"
	urlInput := make(chan Url, BUFFER_SIZE)
	toCrawl := make(chan Url)
	htmlOutput := make(chan Url)
	// extractedLinks := make(chan Url)
	
	// Seeding
	go func() {
		urlInput <- Url{url: homePage}
		// for i := 0; i < 100; i++ {
		// 	urlInput <- Url{url: homePage + "/" + strconv.Itoa(i % 100)}
		// }
	}()

	go urlFilter(urlInput, toCrawl)
	go fetcher(toCrawl, htmlOutput)
	go linkExtractor(htmlOutput, urlInput)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Printf(text)
}
