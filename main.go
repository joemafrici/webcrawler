package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
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
)

// TODO: naked return values only in short functions
func main() {
	baseURL := "https://sive.rs"

	crawl(baseURL)

	/*
		links := links(tree)
		nodes := traverse(tree)
		pullLinks(nodes)
		body, err := io.ReadAll(response.Body)
		checkError(err)

		fmt.Println(string(body))
	*/
}

func crawl(url string) {
	fmt.Printf("crawling url -> %v\n", url)
	response, err := client.Get(url)
	checkError(err)
	defer response.Body.Close()

	tree, err := html.Parse(response.Body)
	checkError(err)

	var urls []string
	extractURLs(tree, &urls)

	for _, url := range urls {
		crawl(url)
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
