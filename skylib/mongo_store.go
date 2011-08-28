package skylib

import (
	"log"
	"os"
	"fmt"
	"flag"
	"launchpad.net/mgo"
	"launchpad.net/gobson/bson"
)


var MongoServer *string = flag.String("mongoServer", "127.0.0.1", "address of mongo server")



func RemoveServiceAt(i int) {

	s := NOS.Services[i]
	c := MC.DB("skynet").C("config")

	err := c.Remove(bson.M{"ipaddress": s.IPAddress, "provides": s.Provides, "port": s.Port,  "protocol": s.Protocol})
	if err != nil {
		log.Panic(err)
	}
	newServices := make([]*RpcService, 0)

	for k, v := range NOS.Services {
		if k != i {
			if v != nil {
				newServices = append(newServices, v)
			}
		}
	}

	NOS.Services = newServices

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
func LoadRegistry() {
	NOS = &RegisteredNetworkServers{}
	NOS.Services = make([]*RpcService, 0)
	var service RpcService
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
		NOS.Services = append(NOS.Services, &service)
	}
	if err != mgo.NotFound {
		log.Panic(err)
	}

}

func (r *RpcService) AddToRegistry() {
	for _, v := range NOS.Services {
		if v != nil {
			if v.Equal(r) {
				LogInfo(fmt.Sprintf("Skipping adding %s : alreday exists.", v.Provides))
				return // it's there so we don't need an update
			}
		}
	}
	NOS.Services = append(NOS.Services, r)

	c := MC.DB("skynet").C("config")

	err := c.Insert(r)
	if err != nil {
		log.Panic(err.String())
	}

}

func (r *RpcService) RemoveFromRegistry() {

	c := MC.DB("skynet").C("config")

	err := c.Remove(bson.M{"ipaddress": r.IPAddress, "provides": r.Provides, "port": r.Port,  "protocol": r.Protocol})
	if err != nil {
		log.Panic(err)
	}

	newServices := make([]*RpcService, 0)

	for _, v := range NOS.Services {
		if v != nil {
			if !v.Equal(r) {
				newServices = append(newServices, v)
			}

		}
	}
	NOS.Services = newServices
}

func RemoveService(i int) {

	s := NOS.Services[i]
	c := MC.DB("skynet").C("config")

	err := c.Remove(bson.M{"ipaddress": s.IPAddress, "provides": s.Provides, "port": s.Port,  "protocol": s.Protocol})
	if err != nil {
		log.Panic(err)
	}
	newServices := make([]*RpcService, 0)

	for k, v := range NOS.Services {
		if k != i {
			if v != nil {
				newServices = append(newServices, v)
			}
		}
	}

	NOS.Services = newServices

}

// Watch for remote changes to the config file.  When new changes occur
// reload our copy of the config file.
// Meant to be run as a goroutine continuously.
func WatchRegistry() {
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


var MC *mgo.Session
//  We *could* use this instead someday.
// var DC *doozer.Conn

// Any Store drop-in file would need to define this global function.
func ConnectStore() {
    MongoConnect()
}
