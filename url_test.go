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

func TestIsSameHost(t *testing.T) {
   test_cases := []struct {
        left, right string
	same bool
   }{
     {"https://www.abc.com", "http://www.abc.com/testpage", true},
     {"https://www.abc.com", "../../testpage", true},
     {"https://www.abc.com", "http://www.abc123.com/testpage", false},
   }

   for i, c := range test_cases {
      got := IsSameHost(c.left, c.right)
      if got != c.same {
          t.Errorf("Case(%d): IsSameHost(%q, %q) == %t, want %t", i, c.left, c.right, got, c.same)
      }
   }
}


func TestResolveUrl(t *testing.T) {
   test_cases := []struct {
        left, right string
	resolved string
   }{
     {"https://www.abc.com", "../../testpage", "https://www.abc.com/testpage"},
     {"https://www.abc.com/dir1/dir2/index.html", "../testpage", "https://www.abc.com/dir1/testpage"},
     {"https://www.abc.com", "http://www.abc123.com/testpage", "http://www.abc123.com/testpage"},
     {"", "http://www.abc123.com/testpage", "http://www.abc123.com/testpage"},
   }

   for i, c := range test_cases {
      got := ResolveUrl(c.left, c.right)
      if got != c.resolved {
          t.Errorf("Case(%d): IsSameHost(%q, %q) == %q, want %q", i, c.left, c.right, got, c.resolved)
      }
   }
}

func TestIsSuccess(t *testing.T) {
   test_cases := []struct {
        url UrlResponse
	success bool
   }{
     {UrlResponse{url: "http://www.abc.com", status: 200}, true},
     {UrlResponse{url: "http://www.abc.com", status: 301}, false},
     {UrlResponse{url: "http://www.abc.com", status: 302, redirectedUrl: "http://www.abc123.com/"}, true},
   }

   for i, c := range test_cases {
      got := c.url.IsSuccess()
      if got != c.success {
          t.Errorf("Case(%d): TestIsSuccess(%q) == %t, want %t", i, c.url, got, c.success)
      }
   }
}
