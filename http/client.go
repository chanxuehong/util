package http

import (
	"net"
	"net/http"
	"time"
)

var DefaultLANClient = http.Client{
	Transport: DefaultLANTransport,
	Timeout:   time.Second * 5,
}

var DefaultWANClient = http.Client{
	Transport: DefaultWANTransport,
	Timeout:   time.Second * 10,
}

var DefaultLANTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          0,
	MaxIdleConnsPerHost:   20,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   2 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var DefaultWANTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   6 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   4 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
