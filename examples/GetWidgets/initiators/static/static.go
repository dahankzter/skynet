package main


import "github.com/bketelsen/skynet/skylib"
import "myStartup"
import "log"
import "os"
import "flag"
import "rpc"

const sName = "Initiator.Static"


// Call the RPC service on the router to process the GetUserDataRequest.
func submitGetUserDataRequest(cr *myStartup.GetUserDataRequest) (*myStartup.GetUserDataResponse, os.Error) {
	var GetUserDataResponse *myStartup.GetUserDataResponse

	client, err := skylib.GetRandomClientByProvides("RouteService.RouteGetUserDataRequest")
	if err != nil {
		if GetUserDataResponse == nil {
			GetUserDataResponse = &myStartup.GetUserDataResponse{}
		}
		GetUserDataResponse.Errors = append(GetUserDataResponse.Errors, err.String())
		return GetUserDataResponse, err
	}
	err = client.Call("RouteService.RouteGetUserDataRequest", cr, &GetUserDataResponse)
	if err != nil {
		if GetUserDataResponse == nil {
			GetUserDataResponse = &myStartup.GetUserDataResponse{}

		}
		GetUserDataResponse.Errors = append(GetUserDataResponse.Errors, err.String())
	}

	return GetUserDataResponse, nil
}


func main() {

	var err os.Error

	// Pull in command line options or defaults if none given
	flag.Parse()

	f, err := os.OpenFile(*skylib.LogFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err == nil {
		defer f.Close()
		log.SetOutput(f)
	}
	skylib.Setup(sName)

	rpc.HandleHTTP()

	for {
		cr := &myStartup.GetUserDataRequest{YourInputValue: "Brian"}
		resp, err := submitGetUserDataRequest(cr)
	if err != nil {
		log.Println(err.String())
	}
		log.Println(resp)
	}


}