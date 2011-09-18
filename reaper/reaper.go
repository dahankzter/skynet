//Copyright (c) 2011 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"flag"
	"github.com/bketelsen/skynet/skylib"
	"fmt"
	"syscall"
	"rpc"
	"rpc/jsonrpc"
	"os"
	"strings"
)

func monitorServices() {
	fmt.Println("Fear the reaper...")

	for {

		skylib.LoadRegistry()
		for _, v := range skylib.NOS.Services {
			var newClient *rpc.Client
			var err os.Error

			hostString := fmt.Sprintf("%s:%d", v.IPAddress, v.Port)
			protocol := strings.ToLower(v.Protocol) // to be safe

			switch protocol {
			default:
				newClient, err = rpc.DialHTTP("tcp", hostString)
			case "json":
				newClient, err = jsonrpc.Dial("tcp", hostString)
			}
			if err != nil {
				//REMOVE HERE
				error := fmt.Sprintf("Bad Service found %s:%d", v.IPAddress, v.Port)
				skylib.LogError(error)
				v.RemoveFromRegistry()
				break
			}
			newClient.Close()

		}
		syscall.Sleep(2e9) // wait 2 seconds
	}
}

func main() {

	flag.Parse()
	skylib.NewAgent().Start()

	// Pull in command line options or defaults if none given
	flag.Parse()

	// normally this happens in Agent, but we're not an agent, so do it ourselves.
	skylib.ConnectStore()

	monitorServices()
}
