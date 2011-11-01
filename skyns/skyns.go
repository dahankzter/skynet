package main

import (
	"github.com/bketelsen/GoServe/serve"
	"os"
	"fmt"
)

type Service struct{
	Name	string
	Address	string
	Timestamp	int64
}

var registry map[string]Service

func cmdPing(*serve.Command) (string, os.Error) {
    return "PONG", nil
}

func cmdRegister(*serve.Command) (string, os.Error) {

    return "REGISTERED", nil
}


func main(){
	var commands = map[string]serve.CmdDescriptor{
	    "PING": {0, cmdPing},
		"REGISTER":{3, cmdRegister},
	}
	
	var cfg = serve.NewConfig("Test server", commands, true)
	
    srv := serve.NewTCPServer(cfg, "127.0.0.1", 7777)
    defer srv.Close()

    go srv.ListenAndServe()
    serve.WaitInt() // Press Ctrl-C to exit

}