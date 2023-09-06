package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/kovalyov-valentin/url-shortener/internal/lib/api/response"
	"github.com/kovalyov-valentin/url-shortener/internal/lib/logger/sl"
	"github.com/kovalyov-valentin/url-shortener/internal/lib/random"
	"github.com/kovalyov-valentin/url-shortener/internal/storage"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: имеет смысл перенести в конфиг или БД
const aliasLenght = 6

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// Функция которая по сути является конструктором для хэндлера
// Во время подключения этого хэндлера к роутеру мы будем возвращать функцию New,
// которая возвращает хэндлер и здесь мы можем передать какие то дополнительные параметры,
// которые будут установлены в каждом обработчике
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Создаем новый валидатор. Говорим что нужно провалидировать структуру req
		if err := validator.New().Struct(req); err != nil {
			// Если найдет ошибку, то вернет ее вот такого типа
			validateErr := err.(validator.ValidationErrors)

			// В чистом виде эту ошибку залогируем без каких либо изменений
			log.Error("invalid request", sl.Err(err))

			//render.JSON(w, r, resp.Error("invalid request"))

			// Формируем запрос в котором сформирована человекочитаемая ошибка
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		// Если alias не назначен, то формируем его из случайных символов
		if alias == "" {
			alias = random.NewRandomString(aliasLenght)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			// Логи storage
			log.Error("failed to add url", sl.Err(err))

			// Логи клиента
			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})

		//responseOK(w, r, alias)
	}
}

//func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
//	render.JSON(w, r, Response{
//		Response: resp.OK(),
//		Alias: alias,
//	})
//}
