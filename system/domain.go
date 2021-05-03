package system

import "net/http"

type API interface {
	version() http.Handler
}

type Service interface {
	version() (string, string, error)
}
