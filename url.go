package main

type Url struct {
	url string
	depth int
	redirectedToUrl string
	html string
	outLinks []string
}

func (url *Url) GetUrl() string {
	if url.redirectedToUrl != "" {
		return url.redirectedToUrl
	} else {
		return url.url
	}
}

func (url *Url) HasOutLinks() bool {
	return len(url.outLinks) > 0
}
