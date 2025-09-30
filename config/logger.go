package config

import (
	"os"
	"time"

	// Add logger
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	// log rotation
	"gopkg.in/natefinch/lumberjack.v2"
	// Add server library
	"net/http"
)

func SetUpLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.Output(&lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10, // МБ
		MaxBackups: 5,
		// MaxAge:     28, // дней
		Compress: true,
	})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func WrapWithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()  // Сохраняем время начала запроса
		next.ServeHTTP(w, r) // Передаём запрос дальше в обработчик
		log.Info().
			Str("method", r.Method).            // HTTP метод (GET, POST)
			Str("url", r.RequestURI).           // URL запроса
			Str("remote", r.RemoteAddr).        // IP адрес клиента
			Dur("duration", time.Since(start)). // Время обработки запроса
			Msg("request completed")            // Сообщение
	})
}