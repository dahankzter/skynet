package main

import (
	"github.com/bketelsen/skynet/skylib"
	"os"
	"flag"
	"time"
	"container/vector"
	"myCompany"
)

var route *skylib.Route

type SubscriptionRouter struct {

}

func callRpcService(service string, operation string, async bool, failOnErr bool, cr myCompany.SubscriptionRequest, rep *myCompany.SubscriptionResponse) (err os.Error) {
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
func (*SubscriptionRouter) RouteSubscriptionRequest(m myCompany.SubscriptionRequest, response *myCompany.SubscriptionResponse) (err os.Error) {
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
	r.Name = "SubscriptionRouter" // I think this doesn't matter any more
	r.LastUpdated = time.Seconds()
	r.Revision = 1
	r.RouteList = new(vector.Vector)

	// Define the chain of services.
	rpcSubscribe := &skylib.RpcCall{Service: "SubscriptionService", Operation: ".ProcessSubscription", Async: false, OkToRetry: false, ErrOnFail: true}

	// Just one, for now.
	r.RouteList.Push(rpcSubscribe)
	/*
		rpcLog := &skylib.RpcCall{Service: "SubscriptionService", Operation: ".LogSubscription", Async: true, OkToRetry: true, ErrOnFail: false}

		// Just one, for now.
		r.RouteList.Push(rpcLog)
	*/
	return
}
