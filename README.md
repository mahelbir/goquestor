# Goquestor

Goquestor is an asynchronous HTTP client based on net/http with additional features.

## Features
- Parallel requests
- Change TLS fingerprint profile
- Change order of TLS extensions
- Change order of headers
- Change order of pseudo headers
- HTTP/1.1 or HTTP/2
- Proxy connection
- Cookie jar

## Installation

Run the following go command
```sh
go get -u github.com/mahelbir/goquestor
```

## Usage

[fhttp Guide](https://github.com/bogdanfinn/fhttp)

[tlsclient Guide](https://bogdanfinn.gitbook.io/open-source-oasis)

```go
package main

import (
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
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
		nil,
		[]tlsclient.HttpClientOption{tlsclient.WithTimeoutSeconds(3)},
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
```

## Credits
Special thanks to the following people for helping me with this project.

[@bogdanfinn](https://github.com/bogdanfinn)

## License

The MIT License (MIT). Please see [License File](LISENCE) for more information.
