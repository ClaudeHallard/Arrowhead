package main

import (
	"net/http"
	"strconv"

	"github.com/hariadivicky/nano"
)

// Todo Controller struct.
type Todo struct{}

// Echo return an OK if the server is running
// route: GET /echo
func (ctrl *Todo) Echo(c *nano.Context) {
	todo := new(Register)

	collection := todo.All()
	c.SetHeader("hello ", "OK")
	c.JSON(http.StatusOK, nano.H{
		"OK":        collection,
		"confirmed": collection,
	})

	//http.StatusOK
}

// route: POST /register
func (ctrl *Todo) Store(c *nano.Context) {
	title := c.PostForm("title")
	isDone, _ := strconv.ParseBool(c.PostForm("is_done"))

	todo := (&Register{
		Title:  title,
		IsDone: isDone,
	}).Save()

	c.JSON(http.StatusOK, nano.H{
		"todo": todo,
	})

}

// query stuff from database.
// route: POST /query
func (ctrl *Todo) Query(c *nano.Context) {

}

// Unregister is functions to delete todo from database.
// route: DELETE /todos/:id
func (ctrl *Todo) Unregister(c *nano.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 8)
	model := new(Register)

	todo := model.Find(uint(id))

	// send http not found when todo does not exists.
	if todo == nil {
		c.String(http.StatusNotFound, "todo not found")
		return
	}

	deleted := todo.Delete()

	c.JSON(http.StatusOK, nano.H{
		"deleted": deleted,
	})
}
