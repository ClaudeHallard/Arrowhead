/*
**********************************************
@authors:
Ivar Grunaeu, ivagru-9@student.ltu.se
Jean-Claude Hallard, jeahal-8@student.ltu.se
Marcus Paulsson, marpau-8@student.ltu.se
Pontus Schünemann, ponsch-9@student.ltu.se
Fabian Widell, fabwid-9@student.ltu.se

**********************************************
starter.go is the compile file that is used to build the service registry.*/

package main

import "serviceRegistry/ServiceRegistry"

func main() {

	ServiceRegistry.StartServiceRegistry()
}
