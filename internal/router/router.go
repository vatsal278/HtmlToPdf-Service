package router

import (
	"github.com/vatsal278/htmltopdfsvc/internal/config"
	"net/http"

	"github.com/PereRohit/util/constant"
	"github.com/PereRohit/util/middleware"
	"github.com/gorilla/mux"

	"github.com/vatsal278/htmltopdfsvc/internal/handler"
)

func Register(container *config.AppContainer) *mux.Router {
	m := mux.NewRouter()

	m.StrictSlash(true)
	m.Use(middleware.RequestHijacker)
	m.Use(middleware.RecoverPanic)
	m = m.PathPrefix("/v1").Subrouter()
	commons := handler.NewCommonSvc()
	m.HandleFunc(constant.HealthRoute, commons.HealthCheck).Methods(http.MethodGet)
	m.NotFoundHandler = http.HandlerFunc(commons.RouteNotFound)
	m.MethodNotAllowedHandler = http.HandlerFunc(commons.MethodNotAllowed)

	// attach routes for services below
	m = attachHtmltopdfsvcRoutes(m, container)

	return m
}

func attachHtmltopdfsvcRoutes(m *mux.Router, c *config.AppContainer) *mux.Router {
	svc := handler.NewHtmltopdfsvc(c)

	m.HandleFunc("/ping", svc.Ping).Methods(http.MethodPost)
	m.HandleFunc("/register", svc.Upload).Methods(http.MethodPost)
	m.HandleFunc("/generate", svc.ConvertToPdf).Methods(http.MethodPost)

	return m
}
