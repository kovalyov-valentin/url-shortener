package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/kovalyov-valentin/url-shortener/internal/lib/api/response"
	"github.com/kovalyov-valentin/url-shortener/internal/lib/logger/sl"
	"github.com/kovalyov-valentin/url-shortener/internal/storage"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("url deleted")

	}

}