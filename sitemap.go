package main

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
