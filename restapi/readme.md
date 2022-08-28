## REST API Exercise
In this exercise, you are tasked with implementing some basic REST endpoints to manage...you guessed it...a todo list!

Specifically, you'll implement:
- `POST /todos`: a route that creates a single new todo
- `GET /todos`: a route that returns the current list of todos
- `GET /todos/{todoID}`: a route that returns a specific todo
- `PUT /todos/{todoID}`: a route that updates a specific todo

Todos consist of a unique ID, a title, and optionally, a user who is assigned...TODO it! :o)

A framework of empty routes is provided in the `rest_api.go` file. All you need to do is fill in the blanks. Note there is no database, no redis, no fanciness whatsoever to this exercise. Anything you need to store can be stored in simple in-memory data structures.

The `rest_api_test.go` file provides a test you can run with `go test` to find any bugs in your implementation. When all the tests pass, you get a gold star :star:!

There is no `main` entrypoint here from which to build a typical binary. Feel free to code one up if you wish, but just passing the provided test suite ought to sufficiently prove your mettle.

Happy coding!
