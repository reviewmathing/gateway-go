package config

import (
	"gateway-go/internal/util"
	"os"
	"path/filepath"
)

const configDirName = "config"

func GetData(fileName string) ([]byte, error) {
	dir, err := util.GetRootDir()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(filepath.Join(dir, configDirName, fileName))
	if err != nil {
		return nil, err
	}
	return file, nil
}
