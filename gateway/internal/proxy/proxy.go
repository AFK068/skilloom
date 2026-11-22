package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// New forwards requests to upstream after removing the gateway route prefix.
func New(upstream, prefix string) http.Handler {
	target, err := url.Parse(upstream)
	switch {
	case err != nil,
		target.Hostname() == "",
		target.User != nil,
		target.Scheme != "http" && target.Scheme != "https":
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "invalid upstream URL", http.StatusBadGateway)
		})
	}

	return http.StripPrefix(prefix, &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.SetURL(target)
			// TODO: Strip Skilloom-specific headers before forwarding.
		},
	})
}
