package main

import (
	"github.com/bketelsen/skynet/skylib"
	"os"
	"flag"
	"message"
)

type StompService struct {

}

func (*StompService) ProcessXmlResponse(m skynetstomp.SkynetStompRequest, response *skynetstomp.SkynetStompResponse) (err os.Error) {

	//Process the message here
	skylib.LogError(m.FirstName)
	//no useful message back to caller
	response.Status = "OK"
	skylib.LogInfo(*response)
	return
}

func main() {
	flag.Parse()
	sig := &StompService{}
	skylib.NewAgent().Register(sig).Start().Wait()
}
