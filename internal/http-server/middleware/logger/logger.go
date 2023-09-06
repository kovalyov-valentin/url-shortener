package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Создаем копию логгера, добавляя подсказку. что это component middleware/logger
		// Это компонент, который будет выводиться с каждой строчкой логов
		log = log.With(
			slog.String("component", "middleware/logger"),
		)
		// Выводим строчку с инфой, чтобы знать что у нас такой хэндлер есть
		log.Info("logger middleware enabled")

		// Содержимое внутренней части хэндлера. Эта часть будет выполняться при каждом запросе
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Эта часть будет выполняться до обработки запроса.
			// То есть пришел запрос, выполняется цепочка хэндлеров middleware
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("reonte)addr", r.RemoteAddr),
				slog.String("user_agent0", r.UserAgent()),
				// Каждому запросу присваивается какой-то request id
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			// Врап здесь нужен, чтобы получить сведения об ответе
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Время, которое ушло на обработку запроса
			t1 := time.Now()

			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			// Передаем управление следующему хэндлеру в цепочке
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
