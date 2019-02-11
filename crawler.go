package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
	"math/rand"
)

func urlFilter(input chan Url, output chan Url) {
	alreadySeen := make(map[string]bool)
	for url := range input {
		resolvedUrl := url.GetUrl()
		if _, isPresent := alreadySeen[resolvedUrl]; !isPresent {
			alreadySeen[resolvedUrl] = true
			output <- url
		}
	}
}

func fetchUrl(url string, out chan struct {x, y string}) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	fmt.Printf("Fetching url %q\n", url)
	out <- struct {x, y string} {url, "html page of url: " + url}
	return
}

func processHtml(url string, html string, out chan string) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	fmt.Printf("Got %d urls from html of %q\n", 1, url)
	out <- url
	return
}

func main() {
	homePage := "www.example.com"
	urlInput := make(chan string)
	htmlOutput := make(chan struct {x, y string})
	extractedLinks := make(chan string)	
	
	// Seeding
	urlInput <- homePage

	go func() {
		for i := 0; i < 100; i++ {
			urlInput <- homePage + "/" + strconv.Itoa(i)
		}
	}()

	for {
		select {
		case url := <-urlInput:
			go fetchUrl(url, htmlOutput)
		case out := <-htmlOutput:
			go processHtml(out.x, out.y, extractedLinks)
		case outLink := <- extractedLinks:
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			fmt.Printf("Outlink %q\n", outLink)
		default:
		}
	}

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Printf(text)
}
