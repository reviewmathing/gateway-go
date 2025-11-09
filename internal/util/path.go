package util

import (
	"os"
	"path/filepath"
	"strings"
)

func GetRootDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return os.Getwd()
	}

	exeDir := filepath.Dir(exePath)

	// 임시 디렉토리 패턴들
	tempPatterns := []string{
		"go-build",         // go run
		"GoLand",           // GoLand
		"Caches/JetBrains", // JetBrains 계열
		"__debug_bin",      // 기타 디버거
		"/tmp/",            // Linux tmp
		"\\Temp\\",         // Windows tmp
	}

	for _, pattern := range tempPatterns {
		if strings.Contains(exePath, pattern) {
			// 임시 경로면 Getwd 사용
			return os.Getwd()
		}
	}

	// 정상 경로면 Executable 사용
	return exeDir, nil
}
