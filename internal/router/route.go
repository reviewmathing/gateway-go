package router

import (
	"fmt"
	"gateway-go/internal/auth"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const RouterConfigName = "config.yml"

const root = "/"
const pathSeparator = "/"

type Router struct {
	routes []Route // 소문자 (외부 노출 불필요)
}

type Route struct {
	Prefix   string `yaml:"prefix"`
	Target   string `yaml:"target"`
	AuthType string `yaml:"auth"`
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
		if !isHTTPScheme(route.Target) {
			return nil, fmt.Errorf("target is not http scheme: target=%q", route.Target)
		}
		if seen[normalize(route.Prefix)] {
			return nil, fmt.Errorf("duplicate route prefix: %q", route.Prefix)
		}
		seen[normalize(route.Prefix)] = true

		authType := route.AuthType
		if authType != "" {
			proxy := auth.Get(authType)
			if proxy == nil {
				return nil, fmt.Errorf("AuthProxy is not setting %s", authType)
			}
		}
	}

	routesCopy := make([]Route, len(config.Routes))
	for i, route := range config.Routes {
		routesCopy[i] = Route{
			Prefix: normalize(route.Prefix),
			Target: normalizeSuffix(route.Target),
		}
	}

	sort.Slice(routesCopy, func(i, j int) bool {
		return len(routesCopy[i].Prefix) > len(routesCopy[j].Prefix)
	})

	return &Router{routes: routesCopy}, nil
}

func (r *Router) Route(path string) (string, string, bool) {
	normalizationPath := normalize(path)
	route, ok := r.matchRoute(normalizationPath)
	if !ok {
		return "", "", false
	}

	target := route.Target
	if route.Prefix == root {
		return target + normalizationPath, route.AuthType, true
	}

	after := normalizationPath[len(route.Prefix):]
	return target + after, route.AuthType, true
}

func (r Router) matchRoute(path string) (Route, bool) {
	for i := range r.routes {
		route := r.routes[i]
		if strings.HasPrefix(path, route.Prefix) {
			remainder := path[len(route.Prefix):]

			if len(remainder) == 0 {
				return route, true
			}

			if strings.HasPrefix(remainder, pathSeparator) {
				return route, true
			}
		}
	}
	return Route{}, false
}

func normalizeSuffix(path string) string {
	if path == root {
		return path
	}
	if path == "" {
		return root
	}
	return strings.TrimSuffix(path, pathSeparator)
}

func normalizePrefix(path string) string {
	if path == "" {
		return root
	}

	if !strings.HasPrefix(path, pathSeparator) {
		return pathSeparator + path
	}
	return path
}

func normalize(path string) string {
	return normalizePrefix(normalizeSuffix(path))
}

func isHTTPScheme(target string) bool {
	return strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://")
}
