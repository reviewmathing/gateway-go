package router

import (
	"gateway-go/internal/logger"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	logger.TestSetUp()

	// 2️⃣ 실제 테스트 실행
	code := m.Run()

	// 프로그램 종료 코드 반환
	os.Exit(code)
}

func TestValidRouter(t *testing.T) {
	yml := `
routes:
  - prefix : /api/test
    target :  http://localhost:8081
  - prefix : /api/test/test
    target : http://localhost:8080
`

	router, err := NewRouter([]byte(yml))
	if err != nil {
		t.Errorf("NewRouter() 반환 에러: %v", err)
	}

	expected := []Route{
		{
			Prefix: "/api/test/test",
			Target: "http://localhost:8080",
		},
		{
			Prefix: "/api/test",
			Target: "http://localhost:8081",
		},
	}

	if !reflect.DeepEqual(router.routes, expected) {
		t.Errorf("라우터가 잘못 생성되었습니다. got=%v, want=%v", router.routes, expected)
	}
}

func TestInvalidRouter(t *testing.T) {
	yml := `
routes:
  - prefix : /api/test
  - prefix : /api/test2
    target : http://localhost:8081
`

	_, err := NewRouter([]byte(yml))
	if err == nil {
		t.Errorf("invalid한 로직 검증 실패")
	}
}

func TestDuplicatePrefixRouter(t *testing.T) {
	yml := `
routes:
  - prefix : /api/test
	target :  http://localhost:8080
  - prefix : /api/test
    target :  http://localhost:8081
`

	_, err := NewRouter([]byte(yml))
	if err == nil {
		t.Errorf("중복 prefix 검증 로직 검증 실패")
	}
}

func TestRoute(t *testing.T) {
	yml := `
routes:
  - prefix : /api/test
    target : http://localhost:8081
  - prefix : /api/test/test
    target : http://localhost:8080
`

	router, _ := NewRouter([]byte(yml))

	targetPath, ok := router.Route("/api/test/test/1")
	if !ok {
		t.Errorf("route실패")
	}

	if targetPath != "http://localhost:8080/1" {
		t.Errorf("Routing 변환 실패: %s", targetPath)
	}
}
