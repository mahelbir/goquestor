package main

import (
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/mahelbir/goquestor"
	"net/url"
)

func main() {
	gq := goquestor.NewGoquestor(200) // concurrency: 200

	// Async multiple requests
	gq.Request(
		http.MethodGet,
		"https://httpbin.org/get",
		nil,
		http.Header{"X-Custom-Header": {"value"}},
		[]tlsclient.HttpClientOption{tlsclient.WithProxyUrl("http://user:pass@127.0.0.1:8080")},
		"request_1_extra_data",
	)
	gq.Request(
		http.MethodPost,
		"https://httpbin.org/post",
		goquestor.EncodeBody(url.Values{"key": {"value"}}),
		http.Header{"content-type": {"application/x-www-form-urlencoded"}},
		nil,
		nil,
	)
	gq.Request(
		http.MethodPost,
		"https://httpbin.org/post",
		goquestor.JSONBody(struct{ Num int }{Num: 1}),
		http.Header{"content-type": {"application/json"}},
		nil,
		nil,
	)
	gq.Request(
		http.MethodPost,
		"https://httpbin.org/delay/5",
		nil,
		nil,
		[]tlsclient.HttpClientOption{tlsclient.WithTimeoutSeconds(3)},
		nil,
	)
	gq.Request(
		http.MethodGet,
		"https://httpbin.org/get?q=1",
		nil,
		http.Header{"user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0"}},
		[]tlsclient.HttpClientOption{tlsclient.WithClientProfile(profiles.Firefox_123)},
		map[string]string{"id": "request_5"},
	)

	// Run all requests
	responses := gq.Execute()

	// Read responses (pointer)
	for _, response := range responses {
		fmt.Print("Status: ")
		fmt.Println(response.Status)
		fmt.Print("Body: ")
		fmt.Println(response.Body)
		fmt.Print("Headers: ")
		fmt.Println(response.Headers)
		fmt.Print("Error: ")
		fmt.Println(response.Error)
		fmt.Print("Identifier: ")
		fmt.Println(*response.Caller)
	}

	// Read responses (slice)
	fmt.Println(goquestor.Result(responses))

}
