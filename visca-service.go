package viscaservice

import (
	"context"
	"net"
)

type EventPublisher interface {
	Publish(context.Context, RequestInfo) error
	Error(context.Context, RequestError) error
}

type RequestInfo struct {
	Action   string
	SourceIP net.IP
	CameraIP net.IP
}

type RequestError struct {
	RequestInfo
	StatusCode int
	Error      error
}
