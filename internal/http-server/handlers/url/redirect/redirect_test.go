package redirect_test

import (
	"github.com/go-chi/chi/v5"
	"github.com/kovalyov-valentin/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/kovalyov-valentin/url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"github.com/kovalyov-valentin/url-shortener/internal/lib/api"
	"github.com/kovalyov-valentin/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(int64(1), tc.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			// Получаем урл, на который произошел редирект
			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			// Check the final URL after redirection
			//require.Equal(t, tc.url, redirectedToURL)
			assert.Equal(t, tc.url, redirectedToURL)

		})
	}
}
