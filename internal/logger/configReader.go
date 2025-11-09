package logger

import (
	"errors"
	"gateway-go/internal/util"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

const LogConfigFileName = "log.yml"

type LogFormat string

const (
	JSON LogFormat = "json"
	TEXT LogFormat = "text"
)

func (lf *LogFormat) parse(name string) error {
	lower := strings.ToLower(name)
	switch lower {
	case "json":
		*lf = JSON
	case "text":
		*lf = TEXT
	default:
		return errors.New("unknown name")
	}
	return nil
}

type logConfig struct {
	App  *ymlLogSetting `yaml:"app"`
	Http *ymlLogSetting `yaml:"http"`
}

func (lc logConfig) appHandler() (slog.Handler, func(), error) {
	var defaultWriter = defaultWrite()
	var appWriter io.WriteCloser
	var writer io.Writer = defaultWriter
	if lc.App != nil {
		appWriter1, err := lc.App.getWriter()
		if err != nil {
			return nil, nil, err
		}
		appWriter = appWriter1

		writer = io.MultiWriter(defaultWriter, appWriter)
		handler, err := toHandler(*lc.App, writer, true)
		if err != nil {
			return nil, nil, err
		}
		return handler, func() {
			_ = appWriter.Close()
			_ = defaultWriter.Close()
		}, nil
	}
	handler, err := toHandler(ymlLogSetting{
		Level:     "INFO",
		LogFormat: "TEXT",
	}, writer, true)
	if err != nil {
		return nil, nil, err
	}

	return handler, func() {
		if appWriter != nil {
			_ = appWriter.Close()
		}
		_ = defaultWriter.Close()
	}, nil
}

func (lc logConfig) httpHandler() (slog.Handler, func(), error) {
	if lc.Http != nil {
		httpWriter, err := lc.Http.getWriter()
		if err != nil {
			return nil, nil, err
		}
		handler, err := toHandler(*lc.Http, httpWriter, false)
		if err != nil {
			return nil, nil, err
		}
		return handler, func() {
			_ = httpWriter.Close()
		}, nil
	}

	return nil, nil, nil
}

func toHandler(yls ymlLogSetting, writer io.Writer, addSource bool) (slog.Handler, error) {
	var logFormat LogFormat
	err := logFormat.parse(yls.LogFormat)
	if err != nil {
		logFormat = TEXT
	}
	var levelTem slog.Level
	err = levelTem.UnmarshalText([]byte(yls.Level))
	if err != nil {
		levelTem = slog.LevelInfo
	}

	option := &slog.HandlerOptions{
		AddSource: addSource,
		Level:     levelTem,
	}

	if logFormat == JSON {
		handler := slog.NewJSONHandler(writer, option)
		return handler, nil
	}
	handler := slog.NewTextHandler(writer, option)
	return handler, nil
}

func defaultWrite() io.WriteCloser {
	return os.Stdout
}

type fileLoggingSetting struct {
	FileName string `yaml:"fileName"`
}

type ymlLogSetting struct {
	Level     string              `yaml:"level"`
	LogFormat string              `yaml:"logFormat"`
	File      *fileLoggingSetting `yaml:"file"`
}

func (yls *ymlLogSetting) getWriter() (io.WriteCloser, error) {
	if yls == nil || yls.File == nil {
		return defaultWrite(), nil
	}
	dir, err := util.GetRootDir()
	if err != nil {
		return nil, errors.New("log path load fail")
	}
	join := filepath.Join(dir, "log/", yls.File.FileName+".log")
	logger := lumberjack.Logger{
		Filename:  join,
		MaxSize:   100,
		MaxAge:    1,
		LocalTime: true,
		Compress:  true,
	}
	return &logger, nil
}

func ReadConfig(data []byte) (Config, error) {
	var config logConfig
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
