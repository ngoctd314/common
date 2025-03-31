package glog

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/ngoctd314/common/env"
)

func SetDefault(prefixEnv string, slogHandler slog.Handler) {
	slog.SetDefault(slog.New(slogHandler))
}

func SlogHandlerWithWriter(prefixEnv string, writer io.Writer, opts *slog.HandlerOptions) slog.Handler {
	switch env.GetString(fmt.Sprintf("%s.format", prefixEnv)) {
	case "json":
		return slog.NewJSONHandler(writer, opts)
	case "text":
		return slog.NewTextHandler(writer, opts)
	default:
		return slog.NewTextHandler(writer, opts)
	}
}

func loggerWriter(prefixEnv string) io.Writer {
	switch env.GetString(fmt.Sprintf("%s.writer", prefixEnv)) {
	case "stdout":
		return os.Stdout
	case "file":
		rotateFile := &lumberjack.Logger{
			Filename:   env.MustString(fmt.Sprintf("%s.file.name", prefixEnv)),
			MaxSize:    env.MustInt(fmt.Sprintf("%s.file.maxSize", prefixEnv)),
			MaxAge:     env.MustInt(fmt.Sprintf("%s.file.maxAge", prefixEnv)),
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   false,
		}

		return rotateFile
	default:
		return os.Stdout
	}
}
