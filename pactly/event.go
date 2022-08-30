package pactly

import (
	"fmt"
	"time"
)

type Header map[string]string

type EventRequest struct {
	Method   string `json:"method" bson:"method"`
	Host     string `json:"host" bson:"host"`
	Path     string `json:"path" bson:"path"`
	Header   Header `json:"header" bson:"header"`
	Body     string `json:"body" bson:"body"`
	Query    string `json:"query" bson:"query"`
	BodySize int    `json:"bodySize" bson:"bodySize"` // contentLength is not reliable, measure body by ourselves
	// todo: should we track basic auth in URL? Might be missing in request header..
}

type EventResponse struct {
	Header     Header `json:"header" bson:"header"`
	Body       string `json:"body" bson:"body"`
	StatusCode int    `json:"statusCode" bson:"statusCode"`
	BodySize   int    `json:"bodySize" bson:"bodySize"` // contentLength is not reliable, measure body by ourselves
}

type Event struct {
	Component       string        `json:"component" bson:"component"`
	Protocol        string        `json:"protocol" bson:"protocol"`
	ProtocolVersion string        `json:"protocolVersion" bson:"protocolVersion"`
	Request         EventRequest  `json:"request" bson:"request"`
	Response        EventResponse `json:"response" bson:"response"`
	Duration        float64       `json:"duration" bson:"duration"`
	Time            time.Time     `json:"time" bson:"time"`
}

func (e *Event) Url() string {
	baseUrl := fmt.Sprintf("%v://%v%v", e.Protocol, e.Request.Host, e.Request.Path)
	if e.Request.Query != "" {
		baseUrl = fmt.Sprintf("%v?%v", baseUrl, e.Request.Query)
	}
	return baseUrl
}
