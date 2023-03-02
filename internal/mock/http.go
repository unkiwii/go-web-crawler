package mock

import (
	"net/http"
	"net/http/httptest"
)

type HTTPGetter func(url string) (*http.Response, error)

func (m HTTPGetter) Get(url string) (*http.Response, error) {
	if m == nil {
		panic("HTTPGetter.Get mock not implemented")
	}
	return m(url)
}

func NewHTTPGetter(statusCode int, response string) HTTPGetter {
	return func(string) (*http.Response, error) {
		w := httptest.NewRecorder()
		w.Write([]byte(response))
		w.WriteHeader(statusCode)
		return w.Result(), nil
	}
}
