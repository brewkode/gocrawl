package main

import (
	"io/ioutil"
	"sync"
	"testing"
)

func TestUrlFilter(t *testing.T) {
	input := make(chan Url, 2)
	output := make(chan Url, 1)
	var wg sync.WaitGroup

	url1 := Url{url: "url1"}
	wg.Add(1)

	// Producer
	go func() {
		defer close(input)
		defer wg.Done()
		input <- url1
		input <- url1
	}()

	// Wait for producer to be done
	wg.Wait()

	// Consumer
	go func() {
		defer close(output)
		urlFilter(input, output)
	}()

	expectedCount := 1
	outputCount := 0

	// Blocking wait on the output channel
	// range over output exits because the consumer does a `defer close(output)`
	// And, this blocking doesn't result in deadlock because these channels are "buffered"
	for x := range output {
		// redundant if-check, as I've not figured out how to drain a channel and count the number of elements without using an intermediate variable - like `x`. and, go fails when `x` is unused
		if &x != nil {
			outputCount++
		}
	}

	if outputCount != expectedCount {
		t.Errorf("UrlFilter did not work. Expected %d, Actual %d", expectedCount, len(output))
	}
}

func TestParse(t *testing.T) {
	test_cases := []struct {
		url, path string
		count     int
	}{
		{"http://brewkode.com", "tests/fixtures/brewkode.html", 6},
	}

	for _, tc := range test_cases {
		outLinks := parse(tc.url, readFile(tc.path, t))

		if len(outLinks) != tc.count {
			t.Errorf("Parse of %q, Outlinks expected %d, actual %d", tc.url, tc.count, len(outLinks))
		}
	}
}

func TestLinkExtractor(t *testing.T) {
	test_cases := []struct {
		url, path string
		count     int
	}{
		{"http://brewkode.com", "tests/fixtures/brewkode.html", 6},
	}

	expected := make(map[string]int)
	for _, tc := range test_cases {
		expected[tc.url] = tc.count
	}

	input := make(chan Url, len(test_cases))
	output := make(chan Url, len(test_cases))
	sitemapInput := make(chan Url, len(test_cases))
	var wg sync.WaitGroup

	wg.Add(1)

	// Producer
	go func() {
		defer close(input)
		defer wg.Done()

		for _, tc := range test_cases {
			input <- Url{url: tc.url, html: readFile(tc.path, t)}
		}
	}()

	// Wait for producer to be done
	wg.Wait()

	// Consumer
	go func() {
		defer close(output)
		linkExtractor(input, output, sitemapInput)
	}()

	// Blocking wait on the output channel
	// range over output exits because the consumer does a `defer close(output)`
	// And, this blocking doesn't result in deadlock because these channels are "buffered"
	for x := range output {
		if len(x.outLinks) != expected[x.url] {
			t.Errorf("ExtractLinks for %q: Expected %d, Actual %d", x.url, expected[x.url], len(x.outLinks))
		}
	}
}

func readFile(path string, t *testing.T) string {
	html, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("Failed reading from file")
	}
	return string(html)
}
