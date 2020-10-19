package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"os"

	// import external libraries
	"github.com/cretz/bine/tor"
	"golang.org/x/net/html"
)

// Usage function
func usage() {
	fmt.Println("Usage: ./simple_to_client http://exempleddfsgsdljfg.onion")
}

func main() {
	// Onion URL must be passed in argument when calling this program
	args := os.Args[1:]
	if len(args) != 1 {
		usage()
		os.Exit(0)
	}
	torAddr := args[0]

	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	fmt.Println("Starting tor and fetching title of ", torAddr,", please wait a few seconds...")
	t, err := tor.Start(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	// Wait at most a minute to start network and get
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()

	// Make connection
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		log.Fatal(err)
	}
	httpClient := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}

	// Create HTTP Get request for Tor Client
//	resp, err := httpClient.Get("https://check.torproject.org")
	resp, err := httpClient.Get(torAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Grab the <title>
	parsed, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Title: %v\n", getTitle(parsed))
}

func getTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		var title bytes.Buffer
		if err := html.Render(&title, n.FirstChild); err != nil {
			panic(err)
		}
		return strings.TrimSpace(title.String())
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := getTitle(c); title != "" {
			return title
		}
	}
	return ""
}
