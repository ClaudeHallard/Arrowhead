/*
**********************************************
@authors:
Ivar Grunaeu, ivagru-9@student.ltu.se
Jean-Claude Hallard, jeahal-8@student.ltu.se
Marcus Paulsson, marpau-8@student.ltu.se
Pontus Schünemann, ponsch-9@student.ltu.se
Fabian Widell, fabwid-9@student.ltu.se

**********************************************
ServiceController.go contains the four service methods that are implemented in this demo
of the service registry, they are called upon when a http request for each of the specific method
is recived by the server listener. */

package ServiceRegistry

import (
	"net/http"
	"strconv"

	"github.com/hariadivicky/nano" // Import RESTful methods.
)

/* Echo return an OK if the server is running
(route: GET serviceregistry/echo) */
func Echo(c *nano.Context) {
	c.JSON(http.StatusOK, "OK")
	println("echo worked")
}

/* Store registers the service and sends back a respons form containing
inserted data to the client.
(route: POST serviceregistry/register) */
func Store(c *nano.Context) {

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

/* Query search for the wanted service and returns a query list containing corresponding services.
(route: POST serviceregistry/query) */
func Query(c *nano.Context) {
	serviceQueryForm := ServiceQueryForm{}
	c.BindJSON(&serviceQueryForm)
	respForm := serviceQueryForm.Query()
	c.JSON(http.StatusOK, respForm)
	println("query worked")
}

/* Unregister delete the specific service that match with the provided data, returns a
confirmation if the service was successfully removed or not.
(route: DELETE serviceregistry/unregister) */
func Unregister(c *nano.Context) {
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
