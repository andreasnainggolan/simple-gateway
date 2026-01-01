package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// New membuat reverse proxy ke target upstream
func New(target string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	// OPTIONAL tapi penting:
	// set ulang Host header agar backend tahu host asli
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = u.Host
	}

	return proxy, nil
}
