package controllers

import (
	"fmt"
	"github.com/IgnasiBosch/go-jwt-auth/api/responses"
	"net/http"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	_ = responses.JSON(w, http.StatusOK, "Welcome home")
}

func (s *Server) PrivateHome(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id")
	fmt.Printf("User id id %s\n", id)

}
