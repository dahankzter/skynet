//Copyright (c) 2011 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package skylib

import (
	"os"
	"fmt"
	"rand"
	"rpc"
	"rpc/jsonrpc"
	"strings"
)


var NOS *RegisteredNetworkServers

// Return a list of all RpcServices which provide the named Service.
func GetAllServiceProviders(classname string) (serverList []*RpcService) {
	fmt.Println("Seeking ", classname)
	fmt.Println("NOS Services Length: ", len(NOS.Services))
	fmt.Println(NOS.Services)

	for _, v := range NOS.Services {
		fmt.Println(v)

		if v != nil && v.Provides == classname {
			serverList = append(serverList, v)
		}
	}
	return
}

func GetAllClientsByService(classname string) (clientList []*rpc.Client) {
	var newClient *rpc.Client
	var err os.Error
	serviceList := GetAllServiceProviders(classname)

	for i, s := range serviceList {
		hostString := fmt.Sprintf("%s:%d", s.IPAddress, s.Port)
		protocol := strings.ToLower(s.Protocol) // to be safe
		switch protocol {
		default:
			newClient, err = rpc.DialHTTP("tcp", hostString)
		case "json":
			newClient, err = jsonrpc.Dial("tcp", hostString)
		}

		if err != nil {
			LogWarn(fmt.Sprintf("Found %d nodes to provide service %s requested on %s, but failed to connect to #%d.",
				len(serviceList), classname, hostString, i))
			//NewError(NO_CLIENT_PROVIDES_SERVICE, classname)
			continue
		}
		clientList = append(clientList, newClient)
	}
	return
}

// This is simple today - it returns the first listed service that matches the request
// Load balancing needs to be applied here somewhere.
func GetRandomClientByService(classname string) (*rpc.Client, os.Error) {
	var newClient *rpc.Client
	var err os.Error
	serviceList := GetAllServiceProviders(classname)

	println("len:", len(serviceList))
	LogWarn("Test LogWarn")
	LogError("Test LogError")
	for len(serviceList) > 0 {
		chosen := rand.Int() % len(serviceList)
		s := serviceList[chosen]

		hostString := fmt.Sprintf("%s:%d", s.IPAddress, s.Port)
		protocol := strings.ToLower(s.Protocol) // to be safe
		switch protocol {
		default:
			newClient, err = rpc.DialHTTP("tcp", hostString)
		case "json":
			newClient, err = jsonrpc.Dial("tcp", hostString)
		}

		if err != nil {
			println(fmt.Sprintf("Found %d nodes to provide service %s requested on %s, but failed to connect.",
				len(serviceList), classname, hostString))
			LogWarn(fmt.Sprintf("Found %d nodes to provide service %s requested on %s, but failed to connect.",
				len(serviceList), classname, hostString))
			s.RemoveFromRegistry()
			l := len(serviceList)
			// We should just remove 'chosen', but for now:
			serviceList = GetAllServiceProviders(classname)
			if l == len(serviceList) {panic("WTF?")}
			continue
			//return nil, NewError(NO_CLIENT_PROVIDES_SERVICE, classname)
		}
		// We have connected.
		return newClient, nil
	}
	println(fmt.Sprintf("Found no node to provide service %s.", classname))
	LogWarn(fmt.Sprintf("Found no node to provide service %s.", classname))
	return nil, NewError(NO_CLIENT_PROVIDES_SERVICE, classname)
}


