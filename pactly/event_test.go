package pactly

import (
	"testing"
)

func TestEvent_Url(t *testing.T) {
	type fields struct {
		Protocol string
		Request  EventRequest
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "test with query", fields: fields{Protocol: "https", Request: EventRequest{Method: "GET", Host: "example.com", Path: "/test", Query: "key=value"}}, want: "https://example.com/test?key=value"},
		{name: "test without query", fields: fields{Protocol: "https", Request: EventRequest{Method: "GET", Host: "example.com", Path: "/test"}}, want: "https://example.com/test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{
				Protocol: tt.fields.Protocol,
				Request:  tt.fields.Request,
			}
			if got := e.Url(); got != tt.want {
				t.Errorf("Url() = %v, want %v", got, tt.want)
			}
		})
	}
}
