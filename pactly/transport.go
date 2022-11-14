package pactly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Transport struct {
	Transport       http.RoundTripper
	serverUrl       *url.URL
	PactlyComponent string
	PactlyToken     string
}

var pendingRequests sync.WaitGroup

// AwaitPendingRequests will wait for all pending pactly requests to finish
func AwaitPendingRequests() {
	pendingRequests.Wait()
}

func DefaultTransport(component string, token string, options *ClientOptions) (*Transport, error) {
	return CustomTransport(options.Url, component, token)
}

func CustomTransport(pactlyUrl string, component string, token string) (*Transport, error) {
	serverUrl, err := url.Parse(pactlyUrl)
	if err != nil {
		return nil, err
	}
	return &Transport{
		Transport:       http.DefaultTransport,
		serverUrl:       serverUrl,
		PactlyComponent: component,
		PactlyToken:     token,
	}, err
}

// RoundTrip is the core part of this module and implements http.RoundTripper.
func (t *Transport) RoundTrip(request *http.Request) (*http.Response, error) {
	ctx := context.Background()
	request = request.WithContext(ctx)

	requestTime := time.Now()
	response, err := t.transport().RoundTrip(request)
	if err != nil {
		return response, err
	}
	responseTime := time.Now()

	var requestBody []byte
	var responseBody []byte

	if request.Body != nil {
		requestBodyRaw, _ := request.GetBody()
		requestBody, err = ioutil.ReadAll(requestBodyRaw)
		if err != nil {
			println(err)
		}
	}
	if response.Body != nil {
		responseBody, _ = ioutil.ReadAll(response.Body)
		err := response.Body.Close()
		if err != nil {
			println(err)
		}
		r := bytes.NewReader(responseBody)
		response.Body = io.NopCloser(r)
	}

	pendingRequests.Add(1)
	go func() {
		defer pendingRequests.Done()
		err := t.logPactlyEvent(*request, requestBody, *response, responseBody, requestTime, responseTime)
		if err != nil {
			fmt.Printf("Failed to log pactly event: %v\n", err)
		}
	}()

	return response, err
}

func (t *Transport) logPactlyEvent(request http.Request, requestBody []byte, response http.Response, responseBody []byte, requestTime time.Time, responseTime time.Time) error {
	if request.URL.Host == t.serverUrl.Host {
		// ignore all requests to pactly for now
		return nil
	}
	requestBodyString := string(requestBody)
	responseBodyString := string(responseBody)

	pactlyEvent := Event{
		Time:            requestTime.UTC(),
		Component:       t.PactlyComponent,
		Protocol:        request.URL.Scheme,
		ProtocolVersion: response.Request.Proto,
		Request: EventRequest{
			Header:   normalizeHeader(request.Header),
			Body:     requestBodyString,
			Host:     request.Host,
			Method:   request.Method,
			Path:     request.URL.Path,
			Query:    request.URL.RawQuery,
			BodySize: len(requestBodyString),
		},
		Response: EventResponse{
			Header:     normalizeHeader(response.Header),
			Body:       responseBodyString,
			BodySize:   len(responseBodyString),
			StatusCode: response.StatusCode,
		},
		Duration: responseTime.Sub(requestTime).Seconds(),
	}

	payload, err := json.Marshal(pactlyEvent)
	if err != nil {
		return err
	}

	resp, err := http.Post(t.serverUrl.String(), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	log.Println("pactly response:" + resp.Status)
	return nil
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func normalizeHeader(header http.Header) Header {
	result := map[string]string{}
	for key, value := range header {
		result[key] = value[0]
	}
	return result
}
