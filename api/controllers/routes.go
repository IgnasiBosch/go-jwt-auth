package controllers

import (
	"github.com/IgnasiBosch/go-jwt-auth/api/middlewares"
	"net/http"
)

func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods(http.MethodGet)
	s.Router.HandleFunc("/home", middlewares.SetMiddlewareJWTAuth(middlewares.SetMiddlewareJSON(s.PrivateHome))).Methods(http.MethodGet)

	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Token)).Methods(http.MethodPost)

}
