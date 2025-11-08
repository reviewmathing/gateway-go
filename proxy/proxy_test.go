package proxy_test

import (
	"fmt"
	"gateway-go/internal/logger"
	"gateway-go/internal/router"
	"gateway-go/proxy"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockRouter struct {
	Routes map[string]string
}

func (m *MockRouter) Route(path string) (string, bool) {
	// 실제 게이트웨이에서는 복잡한 로직이 있겠지만, 테스트를 위해 단순 매핑합니다.
	if target, ok := m.Routes[path]; ok {
		return target, true
	}
	return "", false
}

func TestMain(m *testing.M) {
	logger.TestSetUp()

	// 2️⃣ 실제 테스트 실행
	code := m.Run()

	// 프로그램 종료 코드 반환
	os.Exit(code)
}

func TestProxyHandlerIntegration(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Backend Response OK")
	}))
	defer backend.Close()

	// 2. Mock Router 설정
	mockRouter := &MockRouter{
		Routes: map[string]string{
			// /api/data 경로 요청이 backend 서버의 루트 경로로 프록시되도록 설정 (SetURL 테스트)
			"/api/data": backend.URL + "/api/data",
		},
	}

	proxyHandler := proxy.NewProxy(mockRouter)

	gateway := httptest.NewServer(&proxyHandler)
	defer gateway.Close()

	t.Run("Successful Proxy", func(t *testing.T) {

		// 프록시 서버로 요청 전송
		req, _ := http.NewRequest("GET", gateway.URL+"/api/data", nil)
		// 실제 클라이언트 IP를 시뮬레이션
		req.RemoteAddr = "10.0.0.1:12345"

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatalf("프록시 요청 실패: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("상태 코드 불일치. 기대값: 200, 실제값: %d", resp.StatusCode)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		if string(bodyBytes) != "Backend Response OK" {
			t.Errorf("응답 본문 불일치")
		}
	})

	// --- 테스트 케이스 2: 라우팅 실패 (404 처리) ---
	t.Run("Route Not Found (404)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", gateway.URL+"/api/unknown", nil)
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatalf("프록시 요청 실패: %v", err)
		}
		defer resp.Body.Close()

		// ProxyHandler의 ServeHTTP에서 404를 반환했는지 확인
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("상태 코드 불일치. 기대값: 404, 실제값: %d", resp.StatusCode)
		}
	})
}

func TestProxyHandlerIntegrationWithNotMockRouter(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Backend Response OK")
	}))
	defer backend.Close()

	// 2. Mock Router 설정
	targeturl := backend.URL + "/api/data"
	yamlStr := fmt.Sprintf(`routes:
  - prefix: /api/data
    target: %s`, targeturl)

	newRouter, _ := router.NewRouter([]byte(yamlStr))
	proxyHandler := proxy.NewProxy(newRouter)

	gateway := httptest.NewServer(&proxyHandler)
	defer gateway.Close()

	t.Run("Successful Proxy", func(t *testing.T) {

		// 프록시 서버로 요청 전송
		req, _ := http.NewRequest("GET", gateway.URL+"/api/data?id=10", nil)
		// 실제 클라이언트 IP를 시뮬레이션
		req.RemoteAddr = "10.0.0.1:12345"

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatalf("프록시 요청 실패: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("상태 코드 불일치. 기대값: 200, 실제값: %d", resp.StatusCode)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		if string(bodyBytes) != "Backend Response OK" {
			t.Errorf("응답 본문 불일치")
		}
	})

	// --- 테스트 케이스 2: 라우팅 실패 (404 처리) ---
	t.Run("Route Not Found (404)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", gateway.URL+"/api/unknown", nil)
		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatalf("프록시 요청 실패: %v", err)
		}
		defer resp.Body.Close()

		// ProxyHandler의 ServeHTTP에서 404를 반환했는지 확인
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("상태 코드 불일치. 기대값: 404, 실제값: %d", resp.StatusCode)
		}
	})
}
