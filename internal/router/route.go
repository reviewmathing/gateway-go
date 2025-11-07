package router

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrRouteNotFound = errors.New("route not found")

type Router struct {
	routes []Route // 소문자 (외부 노출 불필요)
}

type Route struct {
	Prefix string `yaml:"prefix"`
	Target string `yaml:"target"`
}

func NewRouter(data []byte) (*Router, error) {
	var config struct {
		Routes []Route `yaml:"routes"`
	}

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	seen := make(map[string]bool)
	for _, route := range config.Routes {
		if route.Prefix == "" || route.Target == "" {
			return nil, fmt.Errorf("invalid route: prefix=%q target=%q",
				route.Prefix, route.Target)
		}
		if !isHTTPSScheme(route.Target) {
			return nil, fmt.Errorf("target is not http scheme: target=%q", route.Target)
		}
		if seen[normalization(route.Prefix)] {
			return nil, fmt.Errorf("duplicate route prefix: %q", route.Prefix)
		}
		seen[normalization(route.Prefix)] = true
	}

	routesCopy := make([]Route, len(config.Routes))
	for i, route := range config.Routes {
		routesCopy[i] = Route{
			Prefix: normalization(route.Prefix),
			Target: normalizationSuffix(route.Target),
		}
	}

	sort.Slice(routesCopy, func(i, j int) bool {
		return len(routesCopy[i].Prefix) > len(routesCopy[j].Prefix)
	})

	return &Router{routes: routesCopy}, nil
}

func (r *Router) Route(path string) (string, error) {
	normalizationPath := normalization(path)
	route, ok := r.matchRoute(normalizationPath)
	if !ok {
		return "", ErrRouteNotFound
	}

	target := route.Target
	if route.Prefix == "/" {
		return target + normalizationPath, nil
	}

	after := normalizationPath[len(route.Prefix):]
	return target + after, nil
}

func (r Router) matchRoute(path string) (Route, bool) {
	for i := range r.routes {
		route := r.routes[i]
		if strings.HasPrefix(path, route.Prefix) {
			remainder := path[len(route.Prefix):]

			if len(remainder) == 0 {
				return route, true
			}

			if strings.HasPrefix(remainder, "/") {
				return route, true
			}
		}
	}
	return Route{}, false
}

func normalizationSuffix(path string) string {
	if path == "/" {
		return path
	}
	if path == "" {
		return "/"
	}
	return strings.TrimSuffix(path, "/")
}

func normalizationPrefix(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func normalization(path string) string {
	return normalizationPrefix(normalizationSuffix(path))
}

func isHTTPSScheme(target string) bool {
	return strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://")
}
