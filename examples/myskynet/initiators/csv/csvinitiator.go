package main

import (
	"net"
	"fmt"
	"myCompany"
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
	service := "SubscriptionRouter"
	client, _ := skylib.GetRandomClientByService(service)

	response := myCompany.SubscriptionResponse{}
	req := myCompany.SubscriptionRequest{}

	req.Name = name
	req.EmailAddress = email

	client.Call(service+".RouteSubscriptionRequest", req, &response)
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
