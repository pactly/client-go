package pactly

const defaultPactlyUrl = "http://localhost:8080/events"

type ClientOptions struct {
	Url string
}

func DefaultOptions() *ClientOptions {
	return &ClientOptions{Url: defaultPactlyUrl}
}

type ClientOption func(o *ClientOptions)

func WithUrl(url string) ClientOption {
	return func(o *ClientOptions) {
		o.Url = url
	}
}
