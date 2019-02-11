package main

import (
       "testing"
       )

func TestGetUrl(t *testing.T) {
   test_cases := []struct {
        url Url
	output string
   }{
     {Url{url: "http://www.abc.com"}, "http://www.abc.com"},
     {Url{url: "http://www.abc.com", redirectedToUrl: "http://www.abc.com/redirectedUrl"}, "http://www.abc.com/redirectedUrl"},
   }

   for i, c := range test_cases {
      got := c.url.GetUrl()
      if got != c.output {
          t.Errorf("Case(%d): GetUrl(%q) == %q, want %q", i, c.url, got, c.output)
      }
   }
}

func TestOutLinks(t *testing.T) {
   test_cases := []struct {
        url Url
	isPresent bool
   }{
     {Url{url: "http://www.abc.com"}, false},
     {Url{url: "http://www.abc.com", redirectedToUrl: "http://www.abc.com/redirectedUrl"}, false},
     {Url{url: "http://www.abc.com", outLinks: append(make([]string, 1), "http://www.abc.com/redirectedUrl")}, true},
   }

   for i, c := range test_cases {
      got := c.url.HasOutLinks()
      if got != c.isPresent {
          t.Errorf("Case(%d): HasOutLinks(%q) == %t, want %t", i, c.url, got, c.isPresent)
      }
   }
}
