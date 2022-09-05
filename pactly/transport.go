package pactly

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Transport struct {
	Transport       http.RoundTripper
	serverUrl       *url.URL
	PactlyComponent string
	PactlyToken     string
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

	print(response.Status)

	// todo async send to pactly in goroutine?
	// ignore all requests to pactly for now
	if request.URL.Host == t.serverUrl.Host {
		return response, nil
	}
	var requestBodyString = ""
	var responseBodyString = ""

	if request.Body != nil {
		requesetBodyRaw, _ := request.GetBody()
		requestBody, err := ioutil.ReadAll(requesetBodyRaw)
		if err != nil {
			print(err)
		}
		requestBodyString = string(requestBody)
	}
	if response.Body != nil {
		responseBody, _ := ioutil.ReadAll(response.Body)
		err := response.Body.Close()
		if err != nil {
			return nil, err
		}
		response.Body = ioutil.NopCloser(bytes.NewBuffer(responseBody))
		responseBodyString = string(responseBody)
	}

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
		log.Fatal(err)
	}

	resp, err := http.Post(t.serverUrl.String(), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("pactly response:" + resp.Status)
	return response, err
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
