package pactly

import "net/http"

// Init starts collecting all http(s) communication and sends them to a pactly server
func Init(component string, token string, options ...ClientOption) error {
	clientOptions := getClientOptions(options)
	pactlyTransport, err := DefaultTransport(component, token, clientOptions)
	if err != nil {
		return err
	}
	http.DefaultTransport = pactlyTransport
	return nil
}

func getClientOptions(options []ClientOption) *ClientOptions {
	clientOptions := DefaultOptions()
	for _, option := range options {
		option(clientOptions)
	}
	return clientOptions
}
