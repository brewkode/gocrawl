package main


import (
	"github.com/parnurzeal/gorequest"
	urllib "net/url"
)


type Url struct {
	url string
	redirectedToUrl string
	html string
	outLinks []string
}

type UrlResponse struct {
	url string
	status int
	redirectedUrl string
	html string
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

func IsSameHost(baseUrl string, extractedUrl string) bool {
	rUrl := ResolveUrl(baseUrl, extractedUrl)
	return Host(baseUrl) == Host(rUrl)
}

func Host(inputUrl string) string {
	u, err := urllib.Parse(inputUrl)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func ResolveUrl(baseUrl string, relativeUrl string) string {
	u, err := urllib.Parse(relativeUrl)
	if err != nil {
		return ""
	}
	b, err := urllib.Parse(baseUrl)
	if err != nil {
		return ""
	}
	return b.ResolveReference(u).String()
}

func NewUrlResponse(url string, resp gorequest.Response, body string) *UrlResponse {
	redirUrl := resp.Request.URL.String() // returns Location header's value if non-empty
	if  redirUrl == url { // case of no redirection
		redirUrl = ""
	}
	return &UrlResponse{url, resp.StatusCode, redirUrl, body}
}

func (resp *UrlResponse) IsSuccess() bool {
	return resp.status == 200 || resp.IsSuccessfulRedirect()
}

func (resp *UrlResponse) IsSuccessfulRedirect() bool {
	return resp.status > 300 && resp.status < 400 && resp.redirectedUrl != ""
}

