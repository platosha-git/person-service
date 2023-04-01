package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"person/user"
)

type Users struct {
	Users []UserOutput `json:"users"`
}

func RegisterHTTPEndpoints(router *mux.Router, uc user.UseCase) {
	handler := NewHandler(uc)

	router.HandleFunc("/api/v1/persons", handler.PersonGetAllHandler).
		Methods(http.MethodOptions, http.MethodGet)
	router.HandleFunc("/api/v1/persons", handler.UserCreateHandler).
		Methods(http.MethodOptions, http.MethodPost)
	router.HandleFunc("/api/v1/persons/{id}", handler.UserGetHandler).
		Methods(http.MethodOptions, http.MethodGet)
	router.HandleFunc("/api/v1/persons/{id}", handler.UserPostHandler).
		Methods(http.MethodOptions, http.MethodPatch)
	router.HandleFunc("/api/v1/persons/{id}", handler.PersonDeleteHandler).
		Methods(http.MethodOptions, http.MethodDelete)
}