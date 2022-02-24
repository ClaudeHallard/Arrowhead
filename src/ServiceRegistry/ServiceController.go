package ServiceRegistry

import (
	"net/http"
	"strconv"

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

	serviceRegistryEntry := ServiceRegistryEntryInput{}
	c.BindJSON(&serviceRegistryEntry)

	println("Registration recived for: " + serviceRegistryEntry.ServiceDefinition)
	respForm := serviceRegistryEntry.Save()

	if respForm == nil {
		println("register denied, service allready exist or bad payload")
		c.JSON(http.StatusOK, "service denied, service allready exist or bad payload")

	} else {
		c.JSON(http.StatusCreated, respForm)
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
func (ctrl *JsonFile) UnregisterOld(c *nano.Context) {
	serviceRegistryEntry := ServiceRegistryEntryInput{}
	c.BindJSON(&serviceRegistryEntry)

	deleted := serviceRegistryEntry.Delete()
	println("Got delete request")
	if deleted {
		c.JSON(http.StatusOK, "Deleted: OK")
	} else {
		c.JSON(http.StatusNotFound, "Deleted FAILED")

	}

}
func (ctrl *JsonFile) Unregister(c *nano.Context) {
	serviceRegistryEntry := ServiceRegistryEntryInput{}
	serviceRegistryEntry.ProviderSystem.Address = c.Query("address")
	serviceRegistryEntry.ProviderSystem.Port, _ = strconv.Atoi(c.Query("port"))
	serviceRegistryEntry.ServiceDefinition = c.Query("service_definition")
	serviceRegistryEntry.ServiceUri = c.Query("service_uri")
	serviceRegistryEntry.ProviderSystem.SystemName = c.Query("system_name")
	deleted := serviceRegistryEntry.Delete()
	println("Got delete request")
	if deleted {
		c.JSON(http.StatusOK, "Deleted: OK")
	} else {
		c.JSON(http.StatusNotFound, "Deleted FAILED")

	}

}
