//Copyright (c) 2012 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"flag"
	"github.com/bketelsen/skynet/skylib"
	"log"
	"time"
)

var (
	BindPort    *int    = flag.Int("port", 9999, "tcp port to listen")
	BindAddr    *string = flag.String("address", "127.0.0.1", "address to bind")
	Region      *string = flag.String("region", "unknown", "region service is located in")
	LogFile     *string = flag.String("logfile", "myservice.log", "name of logfile")
	LogLevel    *int    = flag.Int("loglevel", 1, "log level (1-5)")
  DoozerUri   *string = flag.String("doozer", "", "doozer address")
	DoozerBootUri  *string = flag.String("doozerboot", "", "doozer boot uri")
)

type GetUserDataService struct{}

func (s *GetUserDataService) Registered()   {}
func (s *GetUserDataService) UnRegistered() {}
func (s *GetUserDataService) Started()      {}
func (s *GetUserDataService) Stopped()      {}

const sName = "GetUserData"

type GetUserDataRequest struct {
	YourInputValue string
}

type GetUserDataResponse struct {
	YourOutputValue string
	Errors          []string
}

func NewGetUserDataService() *GetUserDataService {
	r := &GetUserDataService{}
	return r
}

func (ls *GetUserDataService) GetUserData(cr *GetUserDataRequest, lr *GetUserDataResponse) (err error) {
	result := make(chan string)
	timeout := make(chan bool)

	//This function produces the actual result
	go func() {
		time.Sleep(1e8) // force the fail
		result <- " was here"
	}()

	go func() {
		time.Sleep(1e9)
		timeout <- true
	}()

	select {
	case retVal := <-result:
		lr.YourOutputValue = cr.YourInputValue + retVal
	case <-timeout:
		lr.Errors = append(lr.Errors, "Service Timeout")
	}

	return nil
}

func main() {
	// Pull in command line options or defaults if none given
	flag.Parse()

	getDataService := NewGetUserDataService()

	service := skylib.CreateService(getDataService, &skylib.Config{
		Name:                  "GetUserDataService",
		Region:                "Chicago",
		Version:               "1",
		IPAddress:             *BindAddr,
		Port:                  *BindPort,
    DoozerUri:             *DoozerUri,
    DoozerBootUri:         *DoozerBootUri,
	})

	// handle panic so that we remove ourselves from the pool in case of catastrophic failure
	defer func() {
		service.Shutdown()
		if err := recover(); err != nil {
			log.Println("Unrecovered error occured: ", err)
		}
	}()

	// If we pass false here service will not be Registered
	// we could do other work/tasks by implementing the Started method and calling Register() when we're ready
	service.Start(true)
}
