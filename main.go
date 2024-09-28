package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

var collection []string
var deadlinks []string

func unique[T comparable](s []T) []T {
	inResult := make(map[T]bool)
	var result []T
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func getLinks(link string) (links []string) {
	url, err := url.Parse(link)
	if err != nil {
		print(err)
		return links
	}

	// if we don't get a scheme in the links string, we default to https
	if url.Scheme == "" {
		url.Scheme = "https"
	}

	resp, err := http.Get(url.String())
	if err != nil {
		print(err)
		return links
	}

	if resp.StatusCode >= 400 {
		deadlinks = append(deadlinks, url.String())
		return links
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if len(links) > 0 {
		return unique(links)
	}
	return links
}

func main() {
	// if we manage to figure out how to walk the tree downwards, we can just finish this like so
	collection = append(collection, getLinks("robinopletal.com")...)
	collection = append(collection, getLinks("robinopletal.com/posts")...)
	collection = unique(collection)

	for _, v := range collection {
		fmt.Println(v)
	}
}
