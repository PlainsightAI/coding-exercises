package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// this is our in-memory `todo` table
var todos = make([]*todo, 0)

// this is our in-memory `user` table
var users = []*user{
	{ID: 1, Name: "Larry"},
	{ID: 2, Name: "Curly"},
	{ID: 3, Name: "Moe"},
}

// represents a `todo` table record
type todo struct {
	ID         int64
	Title      string
	AssigneeID *int64
}

// represents a `user` table record
type user struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type createTodoInput struct {
	Title      string `json:"title"`
	AssigneeID *int64 `json:"assigneeId"`
}

type updateTodoInput struct {
	Title      string `json:"title"`
	AssigneeID *int64 `json:"assigneeId"`
}

func BuildRestAPI() *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, fmt.Sprintf(`{"message": "unknown route: %s %s"}`, r.Method, r.RequestURI), http.StatusNotFound)
	})

	r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
		var req createTodoInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf(`{"message":%s}`, err.Error()), http.StatusBadRequest)
			return
		}

		// implementation HERE

		http.Error(w, `{"message":"unimplemented"}`, http.StatusInternalServerError)
	})

	r.Get("/todos", func(w http.ResponseWriter, r *http.Request) {
		// implementation HERE

		http.Error(w, `{"message":"unimplemented"}`, http.StatusInternalServerError)
	})

	r.Get("/todos/{todoID}", func(w http.ResponseWriter, r *http.Request) {
		_, err := strconv.Atoi(chi.URLParam(r, "todoID"))
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"message":%s}`, err.Error()), http.StatusBadRequest)
			return
		}

		// implementation HERE

		http.Error(w, `{"message":"unimplemented"}`, http.StatusInternalServerError)
	})

	r.Put("/todos/{todoID}", func(w http.ResponseWriter, r *http.Request) {
		var req updateTodoInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf(`{"message":%s}`, err.Error()), http.StatusBadRequest)
			return
		}

		_, err := strconv.Atoi(chi.URLParam(r, "todoID"))
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"message":%s}`, err.Error()), http.StatusBadRequest)
			return
		}

		// implementation HERE

		http.Error(w, `{"message":"unimplemented"}`, http.StatusInternalServerError)
	})

	return r
}
