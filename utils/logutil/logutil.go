package logutil

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

const LogDir = "/var/log/go-backend"

func Init(appEnv string) {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)

	if appEnv == "production" {
		opts.Level = slog.LevelDebug
		if err := os.MkdirAll(LogDir, 0744); err != nil {
			log.Fatal(err)
		}
		fileName := filepath.Join(LogDir, "go-backend.log")
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}

		handler = slog.NewJSONHandler(file, opts)
	}

	slog.SetDefault(slog.New(handler))
}
