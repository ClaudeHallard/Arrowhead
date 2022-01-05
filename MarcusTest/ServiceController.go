package main

import (
	"net/http"
	"strconv"

	"github.com/hariadivicky/nano"
)

type JsonFile struct{}

// Echo return an OK if the server is running
// route: GET /echo
func (ctrl *JsonFile) Echo(c *nano.Context) {
	ok_mess := "OK"
	println("echo worked")
	c.SetHeader("hello ", "OK")
	c.JSON(http.StatusOK, nano.H{
		"OK": ok_mess,
	})

}

// route: POST /register
func (ctrl *JsonFile) Store(c *nano.Context) {
	ServiceName := c.PostForm("ServiceName")
	metaData := c.PostForm("metaData")
	println("store worked")
	Services := (&Register{
		ServiceName: ServiceName,
		MetaData:    metaData,
	}).Save()

	c.JSON(http.StatusOK, nano.H{
		"Services": Services,
	})

}

// query stuff from database.
// route: POST /query
func (ctrl *JsonFile) Query(c *nano.Context) {
	println("query worked")
	c.SetHeader("query ", "OK")
	c.JSON(http.StatusOK, nano.H{
		"hellohello": "yet not implemented ",
	})
}

// Unregister is functions to delete service from database.

func (ctrl *JsonFile) Unregister(c *nano.Context) {
	id, _ := strconv.ParseUint(c.Param("ID"), 10, 8)
	model := new(Register)
	println("controller", id)
	service := model.Find(uint(id))

	// send http not found when asset does not exists.
	if service == nil {
		c.String(http.StatusNotFound, "query you want to delete does not exist")
		return
	}

	deleted := service.Delete()

	c.JSON(http.StatusOK, nano.H{
		"deleted": deleted,
	})
}
