package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

var (
	config = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	client = &http.Client{
		Transport: transport,
	}

	queue = make(chan string)
    visited = make(map[string]bool)
)

func main() {
	baseURL := "https://js.org"

	go func() {
		queue <- baseURL
	}()

	for href := range queue {
        if !visited[href] {
            crawl(href)
        }
	}
	crawl(baseURL)

}

func fixURL(rawURL string, baseURL string) string {
	url, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
    fmt.Printf("url: %v\n", url) 
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	fixed := base.ResolveReference(url)

	return fixed.String()
}

func crawl(href string) {
	fmt.Printf("crawling url -> %v\n", href)
    visited[href] = true
	response, err := client.Get(href)
	if err != nil {
		fmt.Println(err)
        return
	}
	defer response.Body.Close()

	tree, err := html.Parse(response.Body)
	checkError(err)

	var urls []string
	extractURLs(tree, &urls)

	for _, url := range urls {
		//crawl(fixURL(url, href))
		absoluteURL := fixURL(url, href)
		go func() {
			queue <- absoluteURL
		}()
	}
}

func extractURLs(n *html.Node, urls *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				*urls = append(*urls, attr.Val)
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		extractURLs(child, urls)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ***********************************************
func print(node *html.Node, indent string) {
	fmt.Println(indent + node.Data)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		print(child, indent+" ")
	}
}
func pullLinks(tree []*html.Node) {
	for _, value := range tree {
		if value.Type == html.ElementNode {
			fmt.Println(value.Data)
		}
	}
}

func traverse(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var nodes []*html.Node
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		nodes = append(nodes, traverse(child)...)
	}
	return nodes
}
