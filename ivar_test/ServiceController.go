package main

import (
	"net/http"

	"github.com/hariadivicky/nano"
)

type JsonFile struct{}

// ECHO return an OK if the server is running
func (ctrl *JsonFile) Echo(c *nano.Context) {
	ok_mess := "OK"
	println("echo worked")
	c.SetHeader("hello ", "OK")
	c.JSON(http.StatusOK, nano.H{
		"Echo": ok_mess,
	})

}

// route: POST serviceregistry/register
func (ctrl *JsonFile) Store(c *nano.Context) {

	serviceRegistryEntry := ServiceRegistryEntryInput{}
	c.BindJSON(&serviceRegistryEntry)
	serviceRegistryEntry.Save()
	println("register worked")
	c.JSON(http.StatusOK, nano.H{
		"Services": serviceRegistryEntry,
	})

}

// route: POST serviceregistry/query
func (ctrl *JsonFile) Query(c *nano.Context) {

	serviceQueryForm := ServiceQueryForm{}
	c.BindJSON(&serviceQueryForm)

	serviceQueryList := serviceQueryForm.Query()
	c.JSON(http.StatusOK, nano.H{
		"query": serviceQueryList,
	})
	println("query worked")
}

// route: DELETE serviceregistry/unregister
func (ctrl *JsonFile) Unregister(c *nano.Context) {
	name := c.PostForm("systemName")
	model := new(Register)
	service := model.Find(name)
	if service == nil {
		c.String(http.StatusNotFound, "query you want to delete does not exist")
		return
	}

	deleted := service.Delete(name)
	println("delete worked")
	if deleted {
		c.JSON(http.StatusOK, nano.H{
			"Deleted": "OK",
		})
	}

}
