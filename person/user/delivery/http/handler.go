package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"person/models"
	"person/user"
	"strconv"
)

type Handler struct {
	useCase user.UseCase
}

func NewHandler(useCase user.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type userInput struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

func userInputToModel(user userInput) *models.User {
	return &models.User{
		Name:    user.Name,
		Work:    user.Work,
		Address: user.Address,
		Age:     user.Age,
		Id:      0,
	}
}

func (h *Handler) UserGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.useCase.GetProfile(r.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			msg := map[string]string{"message": "404"}
			err = json.NewEncoder(w).Encode(msg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(UserToUserOutput(u))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PersonGetAllHandler(w http.ResponseWriter, r *http.Request) {
	persons, err := h.useCase.GetAllPersons(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(PersonsToOutputArray(persons))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PersonDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.DeletePerson(r.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			msg := map[string]string{"message": "404"}
			err = json.NewEncoder(w).Encode(msg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UserPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := userInput{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &u); err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(err); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	model := userInputToModel(u)
	newUser, err := h.useCase.ChangeProfile(r.Context(), model, id)

	if err == user.ErrUserNotFound {
		msg := map[string]string{"message": "User not found"}
		err = json.NewEncoder(w).Encode(msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		msg := map[string]string{"message": err.Error()}
		err = json.NewEncoder(w).Encode(msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(UserToUserOutput(newUser)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}

	u := userInput{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &u); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	model := userInputToModel(u)
	userId, err := h.useCase.Create(r.Context(), model)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/v1/persons/%d", userId))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

type UserOutput struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

func UserToUserOutput(user *models.User) *UserOutput {
	return &UserOutput{
		Name:    user.Name,
		Work:    user.Work,
		Address: user.Address,
		Age:     user.Age,
		Id:      user.Id,
	}
}

func PersonsToOutputArray(users []*models.User) []*UserOutput {
	var result []*UserOutput
	for i := 0; i < len(users); i++ {
		result = append(result, UserToUserOutput(users[i]))
	}

	return result
}