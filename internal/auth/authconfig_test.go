package auth

import "testing"

func TestSetUpAuth(t *testing.T) {
	configData := `routes:
  - prefix: /api/test
    target: http://localhost:8081/api/test

auth:
  jwt-auth:
    secret: asdfabsadfasfdfassadfdsafdasfdsafdasfdasfadsfdasfdasf
    auth-header: Authorization
    claims:
      user-id : "userid"
      role: "role"`
	err := SetUpAuth([]byte(configData))
	if err != nil {
		t.Errorf("setup fail %v", err)
	}

	proxy := Get("json")

	if proxy == nil {
		t.Fatal("proxy is not maked")
	}
}
