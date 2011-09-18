//Copyright (c) 2011 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

/*
SkyNet is a framework for a distributed system of processes.

Each process in SkyNet receives its configuration from a centralized configuration repository (currently MongoDB - possibly pluggable in the future).
Configuration changes are polled by each process every few seconds.  This means that starting a new service automatically
advertises that service's availability to the rest of the processes.  This could be improved greatly with mdns or DNS-SD additions.

SkyNet is built on the premise that there will be at least three distinct process types:

Initiators - 
	Initiators are the source of inbound requests.  On a web-centric system, they'd be running HTTP listeners and accept web based requests.

Routers - 
	If Skynet was MVC, then Routers are the "controller" of the system, they call services according to the stored route configuration that matches the request type.
	Technically routers are optional, but if they're not used, Initiators must call Services directly.  In this scenario you lose the capability of changing routes (adding or re-ordering services) in flight.

Services -
	Services are where the work gets done.  These are the processes that service the requests, process the API calls, get the external data, log the requests, authenticate the users, etc.

(Optional) Watchers -
	Watchers are tasks that run and know about the system, but aren't responding to individual requests.  An example of a watcher is the Reaper which pings every Skynet service and removes the unresponsive services from the configuration store.  Another watcher could poll the exported statistics from each node and collect them for publishing in a central web site.

Configuration Storage:
SkyNet uses MongoDB to store configuration data about the available services and routes.  Configuration changes are pushed to MongoDB, and connected clients poll MongoDB periodically to receive changes to the node list.

Transaction Life Cycle:
A typical transaction will come to an Initiator (via http for example) and be sent to a router that is providing the appropriate service to route that type of requests.  The Router checks its routes and calls the services
listed in its route configuration for that Route type.  Routes also define whether a service can be called Asynchronously (fire and forget) or whether the router must wait for a response.  For each service listed in the route
the Router calls the service passing in the request and response objects.  When all services are run, the router returns a response to the Initiator who is responsible for presenting the data to the remote client
appropriately.  In an HTTP example, this could mean translating to data using an HTML template, or an XML template.

TODO:
There are several things that can be improved in SkyNet.  The code needs significant refactoring.  Too much duplication exists.


*/
package skylib
