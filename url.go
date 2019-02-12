package main

import (
	"github.com/parnurzeal/gorequest"	
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

