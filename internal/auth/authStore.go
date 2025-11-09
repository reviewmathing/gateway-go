package auth

import "strings"

var store = map[string]AuthProxy{}

func Save(proxy AuthProxy) {
	store[string(proxy.GetType())] = proxy
}

func Get(typeValue string) AuthProxy {
	return store[strings.ToUpper(typeValue)]
}
