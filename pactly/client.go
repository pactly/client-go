package pactly

import "net/http"

// Init starts collecting all http(s) communication and sends them to a pactly server
func Init(component string, token string) error {
	pactlyTransport, err := DefaultTransport(component, token)
	if err != nil {
		return err
	}
	http.DefaultTransport = pactlyTransport
	return nil
}
