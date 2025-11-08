package proxy

import (
	"context"
	"gateway-go/internal/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type contextKey string

const targetURLKey contextKey = "targetURL"

type Router interface {
	Route(path string) (targetURL string, found bool)
}

type statusCatcherWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusCatcherWriter) WriteHeader(statusCode int) {
	s.status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

type ProxyHandler struct {
	Router Router
	Proxy  httputil.ReverseProxy
}

func NewProxy(router Router) ProxyHandler {
	return ProxyHandler{
		Router: router,
		Proxy: httputil.ReverseProxy{
			Rewrite: func(req *httputil.ProxyRequest) {
				routerDirector(req)
			},
		},
	}
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, ok := p.Router.Route(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		logger.HTTP.LogTransection(*r, http.StatusNotFound)
		return
	}

	ctx := context.WithValue(r.Context(), targetURLKey, path)
	r = r.WithContext(ctx)
	writer := statusCatcherWriter{
		ResponseWriter: w,
		status:         -1,
	}
	p.Proxy.ServeHTTP(&writer, r)
	logger.HTTP.LogTransection(*r, writer.status)
}

func routerDirector(req *httputil.ProxyRequest) {
	path := req.In.Context().Value(targetURLKey).(string)
	targetURL, err := url.Parse(path)
	if err != nil {
		logger.App.Error("url parse miss", "err", err)
		return
	}

	req.Out.Host = targetURL.Host
	req.Out.URL = targetURL
	req.SetXForwarded()
}
