package main

import "strings"

type SiteMap struct {
	outgoing map[string][]string
	incoming map[string][]string
}

type AdjacentUrls struct {
	incoming []string
	outgoing []string
}

func (siteMap *SiteMap) addParent(child string, parent string) {
	if siteMap.incoming == nil {
		siteMap.incoming = make(map[string][]string)
	}

	if elem, ok := siteMap.incoming[child]; ok {
		siteMap.incoming[child] = append(elem, parent)
	} else {
		siteMap.incoming[child] = append([]string{}, parent)
	}
	return
}

func (siteMap *SiteMap) Add(parent string, children ...string) {
	if siteMap.outgoing == nil {
		siteMap.outgoing = make(map[string][]string)
	}
	siteMap.outgoing[parent] = append([]string{}, children...)
	for _, c := range children {
		siteMap.addParent(c, parent)
	}
	return
}

func (siteMap *SiteMap) GetAdjacency(url string) AdjacentUrls {
	adjUrls := AdjacentUrls{}
	if outg, ok := siteMap.outgoing[url]; ok {
		adjUrls.outgoing = outg
	}
	if in, ok := siteMap.incoming[url]; ok {
		adjUrls.incoming = in
	}
	return adjUrls
}

func (siteMap *SiteMap) print(url string, depth int) string {
	visited := make(map[string]bool)
	return build(*siteMap, url, depth, visited) + "\n"
}

func build(siteMap SiteMap, url string, depth int, visited map[string]bool) string {
	entries := []string{}
	entries = append(entries, indented(url, depth+2))
	visited[url] = true
	outLinks, _ := siteMap.outgoing[url]
	for _, outLink := range outLinks {
		if _, present := visited[outLink]; !present {
			entries = append(entries, build(siteMap, outLink, depth+1, visited))
		}
	}
	return strings.Join(entries, "\n")
}

func indented(url string, depth int) string {
	return strings.Repeat(" ", depth) + "- " + url
}
