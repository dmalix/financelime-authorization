package rest

import (
	"github.com/dmalix/authorization-service/app/information"
	"github.com/dmalix/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func Router(logger *zap.Logger, router *mux.Router, handler information.REST, _ middleware.Middleware) {
	router.Handle("/version", handler.Version(logger)).Methods(http.MethodGet)
}
