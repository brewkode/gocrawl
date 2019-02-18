package main

import (
	"testing"
)

func TestGetAdjaceny(t *testing.T) {
	test_cases := []struct {
		url      string
		outgoing []string
		adj      AdjacentUrls
	}{
		{"home", []string{"url1", "url2", "url3"}, AdjacentUrls{outgoing: []string{"url1", "url2", "url3"}}},
		{"url1", []string{"url11", "url21", "home"}, AdjacentUrls{incoming: []string{"home"}, outgoing: []string{"url11", "url21", "url31"}}},
	}
	// intentionally creating a sitemap outside
	// individual test cases will keep adding to same sitemap
	// every test will expect results based on the accumulated graph
	var siteMap SiteMap
	for _, tc := range test_cases {
		siteMap.Add(tc.url, tc.outgoing...)
		output := siteMap.GetAdjacency(tc.url)
		if len(output.incoming) != len(tc.adj.incoming) || len(output.outgoing) != len(tc.adj.outgoing) {
			t.Errorf("GetAdjacency for %q: Expected(%q) != Got(%q)", tc.url, tc.adj, output)
		}
	}
}
