package main


import (
       "testing"
	"io/ioutil"
	"sync"
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
		count int
	}{
		{"http://brewkode.com", "tests/fixtures/brewkode.html", 6},
	}

	for _, tc := range test_cases {
	html, err := ioutil.ReadFile(tc.path)
	if err != nil {
		t.Errorf("Failed reading from file")
	}

	outLinks := parse(tc.url, string(html))
	
	if len(outLinks) != tc.count {
		t.Errorf("Parse of %q, Outlinks expected %d, actual %d", tc.url, tc.count, len(outLinks))
	}
	}
}
