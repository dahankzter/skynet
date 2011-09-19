//Copyright (c) 2011 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

const initiatorTemplate = `package main

import (
	"net"
	"fmt"
	"<%PackageName%>"
	"os"
	"http"
	"bufio"
	"strings"
	"flag"
	"github.com/bketelsen/skynet/skylib"
)

func startHttpServer(addr string) (err os.Error) {
	httpPort, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	go http.Serve(httpPort, nil)
	return
}

func processMessage(name, email string) {
	service := "<%ServiceName%>Router"
	client, _ := skylib.GetRandomClientByService(service)

	response := <%PackageName%>.<%ServiceName%>Response{}
	req := <%PackageName%>.<%ServiceName%>Request{}

	req.Name = name
	req.EmailAddress = email

	client.Call(service+".Route<%ServiceName%>Request", req, &response)
	skylib.Requests.Add(1)
	skylib.LogError("User ", req.Name, " subscribed? ", response.Success)
}

func main() {

	flag.Parse()
	skylib.NewAgent().Start()

	//start management web server
	e := startHttpServer(":8080")
	if e != nil {
		fmt.Println("unable to start http server")
	}

	f, err := os.Open("subscriptions.csv")
	if err != nil {
		panic("Couldn't Open Data File: " + err.String())
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			break
		}
		fields := strings.Split(string(line), ",")
		processMessage(fields[0], fields[1])
	}

}
`
