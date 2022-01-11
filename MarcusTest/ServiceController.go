package main

import (
	"net/http"

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
		"Echo": ok_mess,
	})

}

// route: POST /register
func (ctrl *JsonFile) Store(c *nano.Context) {
	serviceDefinition := c.PostForm("serviceDefinition")
	systemName := c.PostForm("systemName")
	metaData := c.PostForm("metaData")
	port := c.PostForm("Port")
	authenticationInfo := c.PostForm("authenticationInfo")
	serviceURI := c.PostForm("serviceURI")
	endOfValidity := c.PostForm("endOfValidity")
	secure := c.PostForm("secure")
	address := c.PostForm("address")
	version := c.PostForm("version")
	interfaces := c.PostForm("interfaces")

	Services := (&Register{
		ServiceDefinition:  serviceDefinition,
		SystemName:         systemName,
		MetaData:           metaData,
		Port:               port,
		AuthenticationInfo: authenticationInfo,
		ServiceURI:         serviceURI,
		EndOfvalidity:      endOfValidity,
		Secure:             secure,
		Address:            address,
		Version:            version,
		Interfaces:         interfaces,
	}).Save()
	println("register worked")
	c.JSON(http.StatusOK, nano.H{
		"Services": Services,
	})

}

// query stuff from database.
// route: POST /query
func (ctrl *JsonFile) Query(c *nano.Context) {
	name := c.PostForm("systemName")
	model := new(Register)
	service := model.Find(name)

	// send http not found when asset does not exists.
	if service == nil {
		c.String(http.StatusNotFound, "query you are looking for does not exist!")
		return
	}

	queryList := service.Query(name)
	c.JSON(http.StatusOK, nano.H{
		"query": queryList,
	})
	println("query worked")
}

// Unregister is functions to delete service from database.
func (ctrl *JsonFile) Unregister(c *nano.Context) {
	name := c.PostForm("systemName")
	model := new(Register)

	service := model.Find(name)

	// send http not found when asset does not exists.
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
