package restapi_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/sixgill/coding-exercises/restapi"
	"github.com/sixgill/coding-exercises/restapi/util"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TodoWithUser struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Assignee *User  `json:"assignee"`
}

func TestTodos(t *testing.T) {
	r := restapi.BuildRestAPI()
	epTester := util.NewEndpointTester(r)

	newTodo := map[string]interface{}{
		"title":      "My new todo",
		"assigneeId": 2,
	}
	var createResp TodoWithUser
	if err := epTester.Send("POST", "/todos", newTodo, &createResp); err != nil {
		t.Fatal(err)
	}

	if createResp.Title != newTodo["title"] {
		t.Errorf("Expected %s to equal %s", createResp.Title, newTodo["title"])
	}

	t.Run("get", func(t *testing.T) {
		t.Run("GetTodos", testGetTodos(epTester))
		t.Run("GetByID", testGetTodoByID(epTester, createResp.ID))
		t.Run("GetByID404", testGetTodoByID404(epTester, 10))
	})

	t.Run("update", func(t *testing.T) {
		t.Run("UpdateTodo", testUpdateTodo(epTester, createResp.ID))
		t.Run("Update404", testUpdate404(epTester, 10))
	})
}

func testGetTodos(epTester *util.EndpointTester) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		var resp []*TodoWithUser
		if err := epTester.Send("GET", "/todos", nil, &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) != 1 {
			t.Fatal("Expected 1 todo")
		}

		if resp[0].Title != "My new todo" {
			t.Fatalf("Expected %s to equal 'My new todo'", resp[0].Title)
		}

		if resp[0].Assignee == nil {
			t.Error("Expected an assignee for the first todo")
		}

		if resp[0].Assignee.Name != "Curly" {
			t.Errorf("Expected %s to equal %s", resp[0].Assignee.Name, "Curly")
		}
	}
}

func testGetTodoByID(epTester *util.EndpointTester, todoID int64) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		var resp TodoWithUser
		if err := epTester.Send("GET", fmt.Sprintf("/todos/%v", todoID), nil, &resp); err != nil {
			t.Fatal(err)
		}

		if resp.ID != todoID {
			t.Fatalf("Expected the todo ID to be %v, but got %v", todoID, resp.ID)
		}
	}
}

func testGetTodoByID404(epTester *util.EndpointTester, nonExistentTodoID int64) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		var resp TodoWithUser
		err := epTester.Send("GET", fmt.Sprintf("/todos/%v", nonExistentTodoID), nil, &resp)
		if err == nil {
			t.Fatal("Expected an error, but request succeeded")
		}

		if err.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected 404 status code, but got %v", err.StatusCode)
		}
	}
}

func testUpdateTodo(epTester *util.EndpointTester, todoID int64) func(t *testing.T) {
	return func(t *testing.T) {
		update := map[string]interface{}{
			"title":      "My updated todo",
			"assigneeId": 3,
		}
		var resp TodoWithUser
		if err := epTester.Send("PUT", fmt.Sprintf("/todos/%v", todoID), update, &resp); err != nil {
			t.Fatal(err)
		}

		if resp.Title != update["title"] {
			t.Fatalf("Expected %s to equal %s", resp.Title, update["title"])
		}

		if resp.Assignee == nil {
			t.Error("Expected an assignee for the first todo")
		}

		if resp.Assignee.Name != "Moe" {
			t.Errorf("Expected %s to equal %s", resp.Assignee.Name, "Moe")
		}
	}
}

func testUpdate404(epTester *util.EndpointTester, nonExistentTodoID int64) func(t *testing.T) {
	return func(t *testing.T) {
		update := map[string]interface{}{
			"title": "My updated todo",
		}
		var resp TodoWithUser
		err := epTester.Send("PUT", fmt.Sprintf("/todos/%v", nonExistentTodoID), update, &resp)
		if err == nil {
			t.Fatal("Expected an error, but request succeeded")
		}

		if err.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected 404 status code, but got %v", err.StatusCode)
		}
	}
}
