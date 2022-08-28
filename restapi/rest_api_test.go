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

	newTodo1 := map[string]interface{}{
		"title":      "My new todo",
		"assigneeId": 2,
	}
	var createResp1 TodoWithUser
	if err := epTester.Send("POST", "/todos", newTodo1, &createResp1); err != nil {
		t.Fatal(err)
	}

	if createResp1.Title != newTodo1["title"] {
		t.Errorf("Expected %s to equal %s", createResp1.Title, newTodo1["title"])
	}

	newTodo2 := map[string]interface{}{
		"title": "My second todo",
	}
	var createResp2 TodoWithUser
	if err := epTester.Send("POST", "/todos", newTodo2, &createResp2); err != nil {
		t.Fatal(err)
	}

	if createResp2.Title != newTodo2["title"] {
		t.Errorf("Expected %s to equal %s", createResp2.Title, newTodo2["title"])
	}

	t.Run("get", func(t *testing.T) {
		t.Run("GetTodos", testGetTodos(epTester))
		t.Run("GetByID", testGetTodoByID(epTester, createResp1.ID))
		t.Run("GetByID404", testGetTodoByID404(epTester, 10))
	})

	t.Run("update", func(t *testing.T) {
		t.Run("UpdateTodo", testUpdateTodo(epTester, createResp1.ID))
		t.Run("Update404", testUpdate404(epTester, 10))
	})
}

func testGetTodos(epTester *util.EndpointTester) func(t *testing.T) {
	return func(t *testing.T) {
		var resp []*TodoWithUser
		if err := epTester.Send("GET", "/todos", nil, &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) != 2 {
			t.Fatal("Expected 2 todos")
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

		if resp[1].Title != "My second todo" {
			t.Fatalf("Expected %s to equal 'My second todo'", resp[1].Title)
		}

		if resp[1].Assignee != nil {
			t.Fatal("Expected no assignee for second todo")
		}
	}
}

func testGetTodoByID(epTester *util.EndpointTester, todoID int64) func(t *testing.T) {
	return func(t *testing.T) {
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
			t.Error("Expected an assignee")
		}

		if resp.Assignee.Name != "Moe" {
			t.Errorf("Expected %s to equal %s", resp.Assignee.Name, "Moe")
		}

		var getResp []*TodoWithUser
		if err := epTester.Send("GET", "/todos", nil, &getResp); err != nil {
			t.Fatal(err)
		}

		if len(getResp) != 2 {
			t.Fatal("Expected 2 todos")
		}

		if getResp[0].Title != "My updated todo" {
			t.Fatalf("Expected %s to equal 'My updated todo'", getResp[0].Title)
		}

		if getResp[0].Assignee == nil {
			t.Error("Expected an assignee for the first todo")
		}

		if getResp[0].Assignee.Name != "Moe" {
			t.Errorf("Expected %s to equal %s", getResp[0].Assignee.Name, "Moe")
		}

		if getResp[1].Title != "My second todo" {
			t.Fatalf("Expected %s to equal 'My second todo'", getResp[1].Title)
		}

		if getResp[1].Assignee != nil {
			t.Fatal("Expected no assignee for second todo")
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
