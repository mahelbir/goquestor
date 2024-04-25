package goquestor

import (
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tlsclient "github.com/bogdanfinn/tls-client"
	"io"
	"strings"
	"sync"
)

// NewGoquestor Create initializes a Goquestor with a specified concurrency limit.
func NewGoquestor(concurrency int) *Goquestor {
	return &Goquestor{
		clients:     make([]*clientData, 0),
		responses:   make([]*Response, 0),
		concurrency: concurrency,
	}
}

// Request adds a new request to the goquestor along with specific tlsclient options and an identifier/caller.
func (r *Goquestor) Request(method, url string, body io.Reader, headers http.Header, options []tlsclient.HttpClientOption, caller any) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		r.responses = append(r.responses, &Response{
			Error:  fmt.Sprintf("Failed to create HTTP request: %v", err),
			Caller: &caller,
		})
		return
	}
	req.Header = headers

	r.clients = append(r.clients, &clientData{
		request: req,
		options: &options,
		caller:  &caller,
	})
}

// Execute processes all the requests stored in the goquestor up to the concurrency limit.
func (r *Goquestor) Execute() []*Response {
	var wg sync.WaitGroup
	responseChannel := make(chan *Response, len(r.clients))
	semaphore := make(chan struct{}, r.concurrency) // Concurrency control semaphore

	for _, client := range r.clients {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(c *clientData) {
			defer wg.Done()
			defer func() { <-semaphore }()
			respData := doRequest(c.request, c.options, c.caller)
			responseChannel <- respData
		}(client)
	}

	go func() {
		wg.Wait()
		close(responseChannel)
	}()

	responses := make([]*Response, 0)
	for response := range responseChannel {
		responses = append(responses, response)
	}

	r.clients = make([]*clientData, 0)
	r.responses = make([]*Response, 0)

	return responses
}

// doRequest handles the HTTP request execution using a newly created tlsclient and returns the gathered response data.
func doRequest(req *http.Request, options *[]tlsclient.HttpClientOption, caller *any) *Response {
	client, err := tlsclient.NewHttpClient(nil, *options...)
	if err != nil {
		return &Response{
			Error:  fmt.Sprintf("Failed to create HTTP client: %v", err),
			Caller: caller,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return &Response{
			Error:  fmt.Sprintf("HTTP request failed: %v", err),
			Caller: caller,
		}
	}

	responseData := &Response{
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Caller:  caller,
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		responseData.Error = fmt.Sprintf("Failed to read response body: %v", err)
		return responseData
	}
	responseData.Body = string(body)
	responseData.Body = strings.TrimSpace(responseData.Body)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			errorMessage := fmt.Sprintf("Failed to close response body: %v", err)
			if responseData.Error != "" {
				responseData.Error += fmt.Sprintf(" | %s: %v", errorMessage, err)
			} else {
				responseData.Error = fmt.Sprintf("%s: %v", errorMessage, err)
			}
		}
	}()

	return responseData
}

// Result returns the value of responses pointer
func Result(responses []*Response) []Response {
	array := make([]Response, len(responses))
	for _, response := range responses {
		array = append(array, *response)
	}

	return array
}

// Response struct holds all relevant data from a response for later use.
type Response struct {
	Status  int
	Body    string
	Headers http.Header
	Error   string
	Caller  *any
}

// Goquestor struct manages a pool of clients and handles their asynchronous execution.
type Goquestor struct {
	clients     []*clientData
	responses   []*Response
	concurrency int
}

// clientData struct defines the properties of a single request client.
type clientData struct {
	request *http.Request
	options *[]tlsclient.HttpClientOption
	caller  *any
}
