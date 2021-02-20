package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	fmt.Println("Simple http proxy server")
	fmt.Println("------------------------")
	fmt.Println()

	parser := argparse.NewParser("server", "Simple HTTP Proxy server")

	port := parser.Int("p", "port", &argparse.Options{
		Required: false,
		Help:     "Port number to listen to",
		Default:  8080})

	targetUrl := parser.String("t", "target", &argparse.Options{
		Required: false,
		Help:     "Target URL, if is set all request will be forwarded to this base URL"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	fmt.Printf("listening port : %d\n", *port)

	if *targetUrl != "" {
		fmt.Printf("target to use: %s\n", *targetUrl)
	}

	fmt.Println()

	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}

}

func handleRequest(_ http.ResponseWriter, r *http.Request) {

	currentTime := time.Now()

	fmt.Printf("Request recieved %s\n", currentTime.Format("2006-01-02 15:04:05.000000"))
	fmt.Print(formatRequest(r))
	fmt.Printf("\n-------------------- end -------------------------")
}

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			request = append(request, "\n")
			request = append(request, r.Form.Encode())
		} else {
			fmt.Printf("error parsing reqest %v", err)
		}
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
