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
	println("Registration recived for: " + serviceRegistryEntry.ServiceDefinition)
	respForm := serviceRegistryEntry.Save()
	println("register worked")
	c.JSON(http.StatusOK, nano.H{
		"Services": respForm,
	})

}

// route: POST serviceregistry/query
func (ctrl *JsonFile) Query(c *nano.Context) {

	serviceQueryForm := ServiceQueryForm{}
	c.BindJSON(&serviceQueryForm)

	respForm := serviceQueryForm.Query()
	//body, _ := json.Marshal(respForm)
	//c.JSON(respForm)
	c.JSON(http.StatusOK, respForm)
	println("query worked")
}

// route: DELETE serviceregistry/unregister
func (ctrl *JsonFile) Unregister(c *nano.Context) {
	serviceRegistryEntry := ServiceRegistryEntryInput{}
	c.BindJSON(&serviceRegistryEntry)
	println(("test hetr"))
	deleted := serviceRegistryEntry.Delete()
	println("delete worked")
	if deleted {
		c.JSON(http.StatusOK, nano.H{
			"Deleted": "OK",
		})
	} else {
		c.JSON(http.StatusNotFound, nano.H{
			"Deleted": "FAILED",
		})
	}

}
