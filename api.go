package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddr string
	Store      Storage
}

func NewServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		ListenAddr: listenAddr,
		Store:      store,
	}
}

func (server APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/greet", makeHTTPHandlerFunc(greet))
	router.HandleFunc("/anime", makeHTTPHandlerFunc(server.createAnime))
	router.HandleFunc("/anime/{id}", makeHTTPHandlerFunc(server.getAnime))

	http.ListenAndServe(server.ListenAddr, router)
}

func greet(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, map[string]string{"message": "Hello"})
}

func (server APIServer) createAnime(w http.ResponseWriter, r *http.Request) error {
	var request CreateAnimeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()

	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, err)
	}

	anime := Anime{
		Id:     uuid.New().String(),
		Title:  request.Title,
		Author: request.Author,
		Year:   request.Year,
		Status: request.Status,
	}

	err = server.Store.InsertAnime(anime)

	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return writeJSON(w, http.StatusOK, request)
}

func (server APIServer) getAnime(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id := mux.Vars(r)["id"]

		anime, err := server.Store.GetAnime(id)

		if err != nil {
			return writeJSON(w, http.StatusBadRequest, APIError{Error: "id must be an uuid or id does not exist in DB"})
		}

		if anime.Id == "" {
			return writeJSON(w, http.StatusBadRequest, APIError{Error: "id is invalid"})
		}

		return writeJSON(w, http.StatusOK, anime)
	}
	return writeJSON(w, http.StatusBadRequest, APIError{Error: "method not allowed"})
}

type APIError struct {
	Error string `json:"error"`
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
