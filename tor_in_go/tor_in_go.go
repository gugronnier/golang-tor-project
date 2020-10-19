package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"io"

	// import external libraries
	"golang.org/x/net/proxy"
)

func usage() {
	fmt.Println("Usage : ./tor_in_go http://exemple.com")
}

func main() {
	// targeted Onion Site must be enter as argument when calling the program
	args := os.Args[1:]
	if len(args) != 1 {
		usage()
		os.Exit(0)
	}
	torAddr := args[0]

	// Create a socks5 dialer
	dialer, err := proxy.SOCKS5 ("tcp", "127.0.0.1:9050", nil, proxy.Direct)
//	dialer, err := proxy.SOCKS5 ("tcp", "127.0.0.1:9151", nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	// setup HTTP Transport
	tr := &http.Transport{
		Dial: dialer.Dial,
	}
	client := &http.Client{Transport: tr}

//	res, err := client.Get("https://httpbin.org/ip")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	d, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(string(d))


	// Create HTTP request and modify User-Agent
	request, err := http.NewRequest("GET", torAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Linux; rv:59.0) Gecko/59.0 Firefox/59.0")

	// Make request
	response, err := client.Do(request)
        if err != nil {
                log.Fatal(err)
        }

	// Create output file
	outFile, err := os.Create("output.html")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Copy data from HTTP response to ouput file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// show the page after put it in the output file
	page, err := ioutil.ReadFile("output.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(page))

	fmt.Println("###################################################")
	fmt.Println("# This document is also available in `ouput.html` #")
	fmt.Println("###################################################")
}
