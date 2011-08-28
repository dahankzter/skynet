//Copyright (c) 2011 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package skylib

import (
	"log"
	//	"json"
	"flag"
	"os"
	"fmt"
	"rand"
	"rpc"
	"expvar"
	"syscall"
	"os/signal"
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
)

var MC *mgo.Session
var NS *NetworkServers
var RpcServices []*RpcService

var Port *int = flag.Int("port", 9999, "tcp port to listen")
var Name *string = flag.String("name", "changeme", "name of this server")
var BindIP *string = flag.String("bindaddress", "127.0.0.1", "address to bind")
var LogFileName *string = flag.String("logFileName", "myservice.log", "name of logfile")
var LogLevel *int = flag.Int("logLevel", 1, "log level (1-5)")
var MongoServer *string = flag.String("mongoServer", "127.0.0.1", "addr of mongo server")
var Requests *expvar.Int
var Errors *expvar.Int
var Goroutines *expvar.Int
var svc *Service

// This is simple today - it returns the first listed service that matches the request
// Load balancing needs to be applied here somewhere.
func GetRandomClientByProvides(provides string) (*rpc.Client, os.Error) {
	var providesList = make([]*Service, 0)

	var newClient *rpc.Client
	var err os.Error

	for _, v := range NS.Services {
		if v != nil {
			if v.Provides == provides {
				providesList = append(providesList, v)
			}

		}
	}

	if len(providesList) > 0 {
		random := rand.Int() % len(providesList)
		s := providesList[random]

		portString := fmt.Sprintf("%s:%d", s.IPAddress, s.Port)
		newClient, err = rpc.DialHTTP("tcp", portString)
		if err != nil {
			LogWarn(fmt.Sprintf("Found %d Clients to service %s request.", len(providesList), provides))
			return nil, NewError(NO_CLIENT_PROVIDES_SERVICE, provides)
		}

	} else {
		return nil, NewError(NO_CLIENT_PROVIDES_SERVICE, provides)
	}
	return newClient, nil
}

func MongoConnect() {
	var err os.Error
	MC, err = mgo.Mongo("127.0.0.1")
	if err != nil {
		panic(err)
	}

}

// on startup load the configuration file. 
// After the config file is loaded, we set the global config file variable to the
// unmarshaled data, making it useable for all other processes in this app.
func LoadConfig() {
	NS = &NetworkServers{}
	NS.Services = make([]*Service, 0)
	var service Service
	c := MC.DB("skynet").C("config")
	iter, err := c.Find(nil).Iter()
	if err != nil {
		log.Panic(err)
	}
	for {
		err = iter.Next(&service)
		if err != nil {
			break
		}
		fmt.Println(service)
		NS.Services = append(NS.Services, &service)
	}
	if err != mgo.NotFound {
		log.Panic(err)
	}

}

func RemoveServiceAt(i int) {

	s := NS.Services[i]
	c := MC.DB("skynet").C("config")

	err := c.Remove(bson.M{"ipaddress": s.IPAddress, "name": s.Name, "port": s.Port, "provides": s.Provides})
	if err != nil {
		log.Panic(err)
	}
	newServices := make([]*Service, 0)

	for k, v := range NS.Services {
		if k != i {
			if v != nil {
				newServices = append(newServices, v)
			}
		}
	}

	NS.Services = newServices

}

func (r *Service) RemoveFromConfig() {

	c := MC.DB("skynet").C("config")

	err := c.Remove(bson.M{"ipaddress": r.IPAddress, "name": r.Name, "port": r.Port, "provides": r.Provides})
	if err != nil {
		log.Panic(err)
	}

	newServices := make([]*Service, 0)

	for _, v := range NS.Services {
		if v != nil {
			if !v.Equal(r) {
				newServices = append(newServices, v)
			}

		}
	}
	NS.Services = newServices
}

func (r *Service) AddToConfig() {
	for _, v := range NS.Services {
		if v != nil {
			if v.Equal(r) {
				LogInfo(fmt.Sprintf("Skipping adding %s : alreday exists.", v.Name))
				return // it's there so we don't need an update
			}
		}
	}
	NS.Services = append(NS.Services, r)

	c := MC.DB("skynet").C("config")

	err := c.Insert(r)
	if err != nil {
		log.Panic(err.String())
	}

}

// Watch for remote changes to the config file.  When new changes occur
// reload our copy of the config file.
// Meant to be run as a goroutine continuously.
func WatchConfig() {
	/*
		rev, err := DC.Rev()
		if err != nil {
			log.Panic(err.String())
		}
		for {

			// blocking wait call returns on a change
			ev, err := DC.Wait("/servers/config/networkservers.conf", rev)
			if err != nil {
				log.Panic(err.String())
			}
			log.Println("Received new configuration.  Setting local config.")
			setConfig(ev.Body)

			rev = ev.Rev + 1
		}
	*/

}

func initDefaultExpVars(name string) {
	Requests = expvar.NewInt(name + "-processed")
	Errors = expvar.NewInt(name + "-errors")
	Goroutines = expvar.NewInt(name + "-goroutines")
}

func watchSignals() {

	for {
		select {
		case sig := <-signal.Incoming:
			switch sig.(os.UnixSignal) {
			case syscall.SIGUSR1:
				*LogLevel = *LogLevel + 1
				LogError("Loglevel changed to : ", *LogLevel)

			case syscall.SIGUSR2:
				if *LogLevel > 1 {
					*LogLevel = *LogLevel - 1
				}
				LogError("Loglevel changed to : ", *LogLevel)
			case syscall.SIGINT:
				gracefulShutdown()
			}
		}
	}
}

func gracefulShutdown() {
	log.Println("Graceful Shutdown")
	svc.RemoveFromConfig()

	//would prefer to unregister HTTP and RPC handlers
	//need to figure out how to do that
	syscall.Sleep(10e9) // wait 10 seconds for requests to finish  #HACK?
	syscall.Exit(0)
}

func Setup(name string) {
	MongoConnect()
	LoadConfig()
	if x := recover(); x != nil {
		LogWarn("No Configuration File loaded.  Creating One.")
	}

	go watchSignals()

	initDefaultExpVars(name)

	svc = NewService(name)

	svc.AddToConfig()

	go WatchConfig()

	RegisterHeartbeat()

}
