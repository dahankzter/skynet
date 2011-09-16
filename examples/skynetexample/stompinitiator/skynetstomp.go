package main

import (
	"net"
	"fmt"

	stomp "github.com/nf/gostomp"
	"launchpad.net/gobson/bson"
	"message"

	"os"
	"http"
	"flag"
	"github.com/bketelsen/skynet/skylib"
)

var stompport *string = flag.String("stompport", "61613", "stomp port")
var stompserver *string = flag.String("stompserver", "stomp.local", "stomp server")
var login *string = flag.String("login", "skynetuser", "stomp login")
var passcode *string = flag.String("passcode", "changeme", "stomp passcode")

func processStompMessage(request map[string]interface{}) {
	service := "StompRouter"
	client, _ := skylib.GetRandomClientByService(service)

	response := skynetstomp.SkynetStompResponse{}
	req := skynetstomp.SkynetStompRequest{}

	var inquiry map[string]interface{}
	inquiry = request["inquiry"].(bson.M)
	req.FirstName = inquiry["first_name"].(string)
	skylib.LogDebug(req.FirstName)
	client.Call(service+".RouteXmlResponse", req, &response)
	skylib.LogDebug(response.Status)
}

func main() {

	flag.Parse()
	skylib.NewAgent().Start()

	c := connect(*stompserver, *stompport)
	q := "/queue/xml_responses"
	h := stomp.Header{"destination": q}
	_, e := c.Subscribe(h)
	if e != nil {
		fmt.Println("Subscribe error: %v", e)
	}

	for i := range c.Stompdata {
		if i.Error != nil {
			fmt.Println("Receive error: %v", i.Error)
		}

		out := make(map[string]interface{})

		e = bson.Unmarshal(i.Message.Data, out)
		if e != nil {
			fmt.Println("Unmarshal error " + e.String())
		}
		processStompMessage(out)

	}

}

// Connect
func connect(s, p string) (c *stomp.Conn) {
	n, e := net.Dial("tcp", s+":"+p)
	fmt.Println(s, p)
	if e != nil {
		skylib.LogError("Net Connection failed: %v", e)
	}
	h := make(stomp.Header) // empty headers
	h["login"] = *login
	h["passcode"] = *passcode //"horn3tsn3st"
	c, e = stomp.Connect(n, h)
	if e != nil {
		skylib.LogError("CONNECT error 1: %v", e)
	}
	if c == nil {
		skylib.LogError("Connect error 2: %v", c)
	}
	if !c.GetConnected() {
		skylib.LogError("Connection false: %v", c.GetConnected())
	}
	return
}
