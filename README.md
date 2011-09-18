![logo](/bketelsen/skynet/raw/master/documentation/SkyNetLogo.png)


##Introduction
Skynet is a highly available system for building distributed apps in Go.

##Tell me more:
Sknet is a framework for building distributed services in Go. Skynet is built on the premise that services should be composed of one or more building blocks that can be chained together to build a compsite application.

SkyNet is built on the premise that there will be at least three distinct process types:

1. Initiators - Initiators are the source of inbound requests.  On a web-centric system, they'd be running HTTP listeners and accept web based requests.  That isn't required, however.  We have initiators for flat files and TCP connections, too.  If you can get the bytes in using Go, it can be an initiator.
1. Routers - 	Routers are the "controller" of the system, they call services according to the stored route configuration that matches the request type.(Technically routers are optional, but if they're not used, Initiators will call Services directly.  This is an advanced configuration.)
1. Services -Services are where the work gets done.  These are the processes that service the requests, process the API calls, get the external data, log the requests, authenticate the users, etc.  You chain services together in a Route to build an application.
1. (Optional) Watchers -Watchers are tasks that run and know about the system, but aren't responding to individual requests.  An example of a watcher would be a process that watches the other processes in the system and reports on statistics or availability.  The Reaper is a specialized watcher that checks each Skynet cluster member, culling dead processes from the configuration.

##How do I get some of this sweet Skynet action?
Install [Go](http://golang.org) and [MongoDB](http://www.mongodb.org)

	 goinstall github.com/bketelsen/skynet/skygen
	 goinstall github.com/bketelsen/skynet/skylib
	 goinstall github.com/bketelsen/skynet/reaper
	 skygen -packageName=myCompany -serviceName=GetWidgets -targetFullPath="/Users/bketelsen/skynetTest/"


##Do you have pretty diagrams?

Yes.

![logo](/bketelsen/skynet/raw/master/documentation/skynet.jpg)

This rudimentary graph represents a service that is available externally via Web on two different ports, FTP, and TCP.  Each Initiator sends requests to the Router.  The Router looks up the route for this service and calls Logger asynchronously, then Authorization, then GetWidgets.  Adding a new step to the route is as simple as creating the code for the new step (maybe you want to send a notification email, or put some data on a queue, or cache the response) then updating the route.  The Initiators, Routers and other Services don't need to know or care about the new step added to your route.  It all just Works(tm)

##How?
Each process in SkyNet receives its configuration from a centralized configuration repository (currently MongoDB - possibly pluggable in the future).  Configuration changes are pushed to each process when new skynet services are started.  This means that starting a new service automatically
advertises that service's availability to the rest of the members of the skynet cluster.

A typical transaction will come to an Initiator (via http for example) and be sent to a router that is providing the appropriate service to route that type of requests.  The Router checks its routes and calls the services listed in its route configuration for that Route type.  Routes also define whether a service can be called Asynchronously (fire and forget) or whether the router must wait for a response.  For each service listed in the route the Router calls the service passing in the request and response objects.  When all services are run, the router returns a response to the Initiator who is responsible for presenting the data to the remote client appropriately.  In our HTTP example, this could mean translating to data using an HTML template, or an XML/JSON template.

SkyNet uses MongoDB to store configuration data about the available services .  Configuration changes are pushed to Mongo, causing connected clients to immediately become aware of changed configurations.  


##Documentation
There are more (still kind of sparse) documents in the documentation folder.

##Communication
* Group name: Skynet-dev
* Group home page: http://groups.google.com/group/skynet-dev
* Group email address skynet-dev@googlegroups.com

##TODO:
Github Issues now the canonical source of issues for Skynet.

##Open Source - MIT Software License
Copyright (c) 2011 Brian Ketelsen

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
