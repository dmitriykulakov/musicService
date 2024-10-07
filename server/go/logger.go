package server

import (
	"fmt"
	"log/slog"
	"music_service/config"
	"net/http"
	"os"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)
		log := SetupLogger()
		log.Debug(
			fmt.Sprintf("%s %s %s %s",
				r.Method,
				r.RequestURI,
				name,
				time.Since(start)),
		)
	})
}

func SetupLogger() *slog.Logger {
	var log *slog.Logger
	if config.NewLogMode().DebugMode {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
