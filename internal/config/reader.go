package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const configDirName = "config"

func GetData(fileName string) ([]byte, error) {
	path, err := GetRootPath(fileName)
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetRootPath(fileName string) (string, error) {
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		configPath := filepath.Join(exeDir, configDirName, fileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}

	wd, err := os.Getwd()
	if err == nil {
		configPath := filepath.Join(wd, configDirName, fileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}
	}
	return "", fmt.Errorf(fileName, " not found in any known location")
}
