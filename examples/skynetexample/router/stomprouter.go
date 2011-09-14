package main

import (
	"github.com/bketelsen/skynet/skylib"
	"os"
	"flag"
	"time"
	"message"
	"container/vector"
)

var route *skylib.Route

type StompRouter struct {

}

func callRpcService(service string, operation string, async bool, failOnErr bool, cr skynetstomp.SkynetStompRequest, rep *skynetstomp.SkynetStompResponse) (err os.Error) {
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
func (*StompRouter) RouteXmlResponse(m skynetstomp.SkynetStompRequest, response *skynetstomp.SkynetStompResponse) (err os.Error) {
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
	sig := &StompRouter{}
	skylib.NewAgent().Register(sig).Start().Wait()
}

// To change a route, make a new router and kill
// the old ones.
func CreateInitialRoute() (r *skylib.Route) {

	r = new(skylib.Route)
	// Create a basic Route object.
	r.Name = "StompRouter" // I think this doesn't matter any more
	r.LastUpdated = time.Seconds()
	r.Revision = 1
	r.RouteList = new(vector.Vector)

	// Define the chain of services.
	rpcScore := &skylib.RpcCall{Service: "StompService", Operation: ".ProcessXmlResponse", Async: false, OkToRetry: false, ErrOnFail: true}

	// Just one, for now.
	r.RouteList.Push(rpcScore)
	return
}
