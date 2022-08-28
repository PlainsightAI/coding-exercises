package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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

type hydratedTodo struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Assignee *user  `json:"assignee"`
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

		newTodo := todo{
			ID:         time.Now().UnixNano(),
			Title:      req.Title,
			AssigneeID: req.AssigneeID,
		}

		todos = append(todos, &newTodo)

		// for now the tests don't assert that this returns an assignee-hydrated todo
		render.Respond(w, r, todoWithAssignee(&newTodo))
	})

	r.Get("/todos", func(w http.ResponseWriter, r *http.Request) {
		// the join here offers a nice chance to ask about the naive quadratic double loop vs linear-ish time implementation you get with a hash map
		usersByID := mapUsersByID()
		resp := make([]*hydratedTodo, len(todos))

		for i, todo := range todos {
			var assignee *user
			if todo.AssigneeID != nil {
				assignee = usersByID[*todo.AssigneeID]
			}

			resp[i] = &hydratedTodo{
				ID:       todo.ID,
				Title:    todo.Title,
				Assignee: assignee,
			}
		}

		render.Respond(w, r, resp)
	})

	r.Get("/todos/{todoID}", func(w http.ResponseWriter, r *http.Request) {
		todoID, err := strconv.Atoi(chi.URLParam(r, "todoID"))
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"message":%s}`, err.Error()), http.StatusBadRequest)
			return
		}

		for _, todo := range todos {
			if todo.ID == int64(todoID) {
				render.Respond(w, r, todoWithAssignee(todo))
				return
			}
		}

		http.Error(w, fmt.Sprintf(`{"message":"no todo found for id: %v"}`, todoID), http.StatusNotFound)
	})

	r.Put("/todos/{todoID}", func(w http.ResponseWriter, r *http.Request) {
		var req updateTodoInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf(`{"message":%s}`, err.Error()), http.StatusBadRequest)
			return
		}

		todoID, err := strconv.Atoi(chi.URLParam(r, "todoID"))
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"message":%s}`, err.Error()), http.StatusBadRequest)
			return
		}

		updatedTodoIdx := -1

		for i, todo := range todos {
			if todo.ID == int64(todoID) {
				updatedTodoIdx = i
				break
			}
		}

		if updatedTodoIdx == -1 {
			http.Error(w, fmt.Sprintf(`{"message":"no todo found for id: %v"}`, todoID), http.StatusNotFound)
			return
		}

		updatedTodo := todos[updatedTodoIdx]
		updatedTodo.Title = req.Title
		updatedTodo.AssigneeID = req.AssigneeID

		// the tests DO assert that this one is assignee-hydrated
		render.Respond(w, r, todoWithAssignee(updatedTodo))
	})

	return r
}

func mapUsersByID() map[int64]*user {
	usersByID := make(map[int64]*user)

	for _, user := range users {
		usersByID[user.ID] = user
	}

	return usersByID
}

func todoWithAssignee(todo *todo) *hydratedTodo {
	var assignee *user
	if todo.AssigneeID != nil {
		for _, user := range users {
			if *todo.AssigneeID == user.ID {
				assignee = user
				break
			}
		}
	}

	return &hydratedTodo{
		ID:       todo.ID,
		Title:    todo.Title,
		Assignee: assignee,
	}
}
