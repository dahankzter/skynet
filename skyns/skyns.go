package main

import (
	"github.com/bketelsen/GoServe/serve"
	"os"
	"time"
	"fmt"

)

type ServiceArray []*Service

type Service struct{
	Name	string
	Address	string
	Timestamp	int64
}

var registry map[string]ServiceArray

func cmdPing(*serve.Command) (string, os.Error) {
    return "PONG", nil
}

func cmdRegister(cmd *serve.Command) (string, os.Error) {
	service := &Service{Name: cmd.Args[0], Address: cmd.Args[1], Timestamp: time.Seconds()}
	registry[service.Name] = registry[service.Name].replaceOrAdd(service)
    return "REGISTERED " + cmd.Args[0],  nil
}

func cmdGet(cmd *serve.Command) (string, os.Error) {
	 sa := registry[cmd.Args[0]] 
	fmt.Println(cmd.Args[0])
	fmt.Println("Found sa:", sa, len(sa))
	if len(sa) > 0 {
		return "FOUND " + sa[0].Name + " " + sa[0].Address, nil
	}
    return "NOTFOUND" , os.NewError("Not Found")
}

func (this *Service) equals(that *Service) bool{
	return (this.Name == that.Name) && (this.Address == that.Address)
}

func (r ServiceArray) replaceOrAdd(s *Service)(ServiceArray) {
	var x int
	for _, item := range r{
		x = x+1
		if item.equals(s){
			item.Timestamp = s.Timestamp
			return r;
		}
	}
	//fall through to add
	n := append(r,s)
	fmt.Println(n)
	return n
}

func main(){

	registry = make(map[string]ServiceArray)
	
	var commands = map[string]serve.CmdDescriptor{
	    "PING": {0, cmdPing},
		"REGISTER":{2, cmdRegister},
		"GET":{1,cmdGet},
	}
	
	var cfg = serve.NewConfig("Skynet Registry Service", commands, true)
	
    srv := serve.NewTCPServer(cfg, "127.0.0.1", 7777)
    defer srv.Close()

    go srv.ListenAndServe()
    serve.WaitInt() // Press Ctrl-C to exit

}