//Copyright (c) 2011 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

const routerTemplate = `package main

import (
	"github.com/bketelsen/skynet/skylib"
	"os"
	"flag"
	"time"
	"container/vector"
	"<%PackageName%>"
)
	

var route *skylib.Route

type <%ServiceName%>Router struct {

}

func callRpcService(service string, operation string, async bool, failOnErr bool, cr <%PackageName%>.<%ServiceName%>Request, rep *<%PackageName%>.<%ServiceName%>Response) (err os.Error) {
	defer skylib.CheckError(&err)

	rpcClient, err := skylib.GetRandomClientByService(service)
	if err != nil {
		skylib.LogError("No service provides", service)
		if failOnErr {
			return skylib.NewError(skylib.NO_CLIENT_PROVIDES_SERVICE, service)
		} else {
			return nil
		}
	}
	name := service + operation
	if async {
		go rpcClient.Call(name, cr, rep)
		skylib.LogInfo("Called service async", name)
		return nil
	}
	skylib.LogInfo("Calling : " + name)
	err = rpcClient.Call(name, cr, rep)
	if err != nil {
		skylib.LogError("failed connection, retrying", err)
		// get another one and try again!
		rpcClient, err := skylib.GetRandomClientByService(service)
		err = rpcClient.Call(name, cr, rep)
		if err != nil {
			return skylib.NewError(err.String(), name)
		}
	}
	skylib.LogInfo("Called service operation sync", name)
	return nil
}


// Service operation for RPC.
func (*<%ServiceName%>Router) Route<%ServiceName%>Request(m <%PackageName%>.<%ServiceName%>Request, response *<%PackageName%>.<%ServiceName%>Response) (err os.Error) {
	defer skylib.CheckError(&err)
	skylib.LogInfo(route)
	for i := 0; i < route.RouteList.Len(); i++ {
		rpcCall := route.RouteList.At(i).(*skylib.RpcCall)

		err := callRpcService(rpcCall.Service, rpcCall.Operation, rpcCall.Async, rpcCall.ErrOnFail, m, response)
		if err != nil {
			skylib.Errors.Add(1)
			return err
		}

	}

	skylib.Requests.Add(1)
	return nil

}


func main() {
	flag.Parse()
	route = CreateInitialRoute()
	sig := &SubscriptionRouter{}
	skylib.NewAgent().Register(sig).Start().Wait()
}


// checkError is a deferred function to turn a panic with type *Error into a plain error return.
// Other panics are unexpected and so are re-enabled.
func checkError(error *os.Error) {
	if v := recover(); v != nil {
		if e, ok := v.(*skylib.Error); ok {
			*error = e
		} else {
			// runtime errors should crash
			panic(v)
		}
	}
}

// To change a route, make a new router and kill
// the old ones.
func CreateInitialRoute() (r *skylib.Route) {

	r = new(skylib.Route)
	// Create a basic Route object.
	r.Name = "<%ServiceName%>Router" // I think this doesn't matter any more
	r.LastUpdated = time.Seconds()
	r.Revision = 1
	r.RouteList = new(vector.Vector)

	// Define the chain of services.
	rpcScore := &skylib.RpcCall{Service: "<%ServiceName%>Service", Operation: ".Process<%ServiceName%>", Async: false, OkToRetry: false, ErrOnFail: true}

	// Just one, for now.
	r.RouteList.Push(rpcScore)
	return
}
`
