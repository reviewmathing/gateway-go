package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

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
				path, ok := router.Route(req.In.URL.Path)

				if !ok {
					return
				}

				routerDirector(path, req)
			},
		},
	}
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, ok := p.Router.Route(r.URL.Path)
	if !ok {
		http.NotFound(w, r) // 라우터에 없으면 404
		return
	}

	p.Proxy.ServeHTTP(w, r) // Proxy 호출
}

func routerDirector(path string, req *httputil.ProxyRequest) {
	targetURL, err := url.Parse(path)
	if err != nil {
		return
	}

	req.Out.URL = targetURL
	req.SetXForwarded()
}
