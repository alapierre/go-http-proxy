package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var headersToSkip = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

var skipSet map[string]bool

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

	skipSet = prepareHeadersSet(headersToSkip)

	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}

}

func handleRequest(resp http.ResponseWriter, req *http.Request) {

	currentTime := time.Now()

	fmt.Printf("Request recieved %s\n", currentTime.Format("2006-01-02 15:04:05.000000"))
	fmt.Print(formatRequest(req))
	fmt.Printf("\n-------------------- request end -------------------------\n")

	//client := &http.Client{}
	//req.RequestURI = ""

	//cleanHeaders(req.Header)
	//
	//if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
	//	appendXForwardHeader(req.Header, clientIP)
	//}
	//
	//tmp, err := client.Do(req)
	//if err != nil {
	//	http.Error(resp, "Server Error", http.StatusInternalServerError)
	//	log.Fatal("ServeHTTP:", err)
	//}
	//defer closeSilently(tmp.Body)
	//
	//fmt.Printf("%s status: %v", req.RemoteAddr, tmp.Status)
	//
	//copyHeaders(tmp.Header, resp.Header(), skipSet)
	//resp.WriteHeader(tmp.StatusCode)
	//if _, err := io.Copy(resp, tmp.Body); err != nil {
	//	fmt.Println("Error coping response from client")
	//}
	//
	//fmt.Print(formatResponse(tmp))
	//fmt.Printf("-------------------- responce end -------------------------")

}

func formatResponse(r *http.Response) string {

	var res []string

	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			res = append(res, fmt.Sprintf("%v: %v", name, h))
		}
	}

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		res = append(res, string(body))
	} else {
		fmt.Printf("Error getting responce body")
	}

	return strings.Join(res, "\n")
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

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("=== body START ===")
		fmt.Println(string(bytes))
		fmt.Println("=== body END ===")

		//if err := r.ParseForm(); err != nil {
		//	request = append(request, "\n")
		//	request = append(request, r.Form.Encode())
		//} else {
		//	fmt.Printf("error parsing reqest %v", err)
		//}
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func prepareHeadersSet(headers []string) map[string]bool {
	res := make(map[string]bool)
	for _, header := range headers {
		res[header] = true
	}
	return res
}

func copyHeaders(src, dest http.Header, toSkip map[string]bool) {
	for name, header := range src {
		if !toSkip[name] {
			for _, v := range header {
				dest.Add(name, v)
			}
		}
	}
}

func appendXForwardHeader(header http.Header, host string) {
	if prior, ok := header["X-Forwarded-For"]; ok {
		prior = append(prior, host)
		host = strings.Join(prior, ", ")
	}
	header.Set("X-Forwarded-For", host)
}

func checkSchema(resp http.ResponseWriter, req *http.Request) bool {

	if req.URL.Scheme == "http" || req.URL.Scheme == "https" {
		return true
	}

	msg := "unsupported protocal scheme " + req.URL.Scheme
	http.Error(resp, msg, http.StatusBadRequest)
	log.Println(msg)
	return false
}

func cleanHeaders(header http.Header) {
	for _, h := range headersToSkip {
		header.Del(h)
	}
}

func closeSilently(c io.Closer) {
	_ = c.Close()
}
