package proxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type contextKey string

const targetURLKey contextKey = "targetURL"

type Router interface {
	Route(path string) (targetURL string, found bool)
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
		http.NotFound(w, r) // 라우터에 없으면 404
		return
	}

	ctx := context.WithValue(r.Context(), targetURLKey, path)
	r = r.WithContext(ctx)
	p.Proxy.ServeHTTP(w, r)
}

func routerDirector(req *httputil.ProxyRequest) {
	path := req.In.Context().Value(targetURLKey).(string)
	targetURL, err := url.Parse(path)
	if err != nil {
		return
	}

	req.Out.URL = targetURL
	req.SetXForwarded()
}
