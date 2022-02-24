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

	c.JSON(http.StatusOK, ok_mess)

}

// route: POST serviceregistry/register
func (ctrl *JsonFile) Store(c *nano.Context) {

	//bind
	ServiceRegistryEntry := ServiceRegistryEntryInput{}
	c.BindJSON(&ServiceRegistryEntry)

	// If there is any metadata in struct --> run saveJava
	if ServiceRegistryEntry.MetadataJava.AdditionalProp1 != "" {
		respform1 := ServiceRegistryEntry.SaveJava() // måste kolla struct istället
		println(ServiceRegistryEntry.MetadataJava.AdditionalProp1)
		if respform1 == nil {
			println("register denied, service allready exist")
			println("java")
			c.JSON(http.StatusOK, "service denied")
		} else {
			println("register Java")
			println("register worked")
			c.JSON(http.StatusOK, respform1)
		}
	} else {

		respform2 := ServiceRegistryEntry.Save()
		if respform2 == nil {
			println("golang")
			println("register denied, service allready exist")
			c.JSON(http.StatusOK, "service denied")
		} else {
			println("register golang")
			println("register worked")
			c.JSON(http.StatusOK, respform2)
		}
	}
}

// route: POST serviceregistry/query
func (ctrl *JsonFile) Query(c *nano.Context) {

	serviceQueryForm := ServiceQueryForm{}
	c.BindJSON(&serviceQueryForm)

	respForm := serviceQueryForm.Query()
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
		c.JSON(http.StatusOK, "Deleted: OK")
	} else {
		c.JSON(http.StatusNotFound, "Deleted FAILED")

	}

}
