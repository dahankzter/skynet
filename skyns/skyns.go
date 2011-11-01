package main

import (
	"github.com/bketelsen/GoServe/serve"
	"os"
	"time"
)

type Service struct{
	Name	string
	Address	string
	Timestamp	int64
}

var registry map[string]*Service

func cmdPing(*serve.Command) (string, os.Error) {
    return "PONG", nil
}

func cmdRegister(cmd *serve.Command) (string, os.Error) {
	service := &Service{Name: cmd.Args[0], Address: cmd.Args[1], Timestamp: time.Seconds()}
	registry[service.Name] = service
    return "REGISTERED " + cmd.Args[0],  nil
}

func cmdGet(cmd *serve.Command) (string, os.Error) {
	 service := registry[cmd.Args[0]] 
    return "FOUND " + service.Name + " " + service.Address,  nil
}



func main(){
	
	registry = make(map[string]*Service)
	
	var commands = map[string]serve.CmdDescriptor{
	    "PING": {0, cmdPing},
		"REGISTER":{2, cmdRegister},
		"GET":{1,cmdGet},
	}
	
	var cfg = serve.NewConfig("Test server", commands, true)
	
    srv := serve.NewTCPServer(cfg, "127.0.0.1", 7777)
    defer srv.Close()

    go srv.ListenAndServe()
    serve.WaitInt() // Press Ctrl-C to exit

}