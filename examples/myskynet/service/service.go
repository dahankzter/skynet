package main

import (
	"github.com/bketelsen/skynet/skylib"
	"os"
	"flag"
	"myCompany"
)

type SubscriptionService struct {

}

func (*SubscriptionService) LogSubscription(m myCompany.SubscriptionRequest, response *myCompany.SubscriptionResponse) (err os.Error) {

	//Log the message here

	skylib.LogError(m)
	skylib.Requests.Add(1)
	return
}

func (*SubscriptionService) ProcessSubscription(m myCompany.SubscriptionRequest, response *myCompany.SubscriptionResponse) (err os.Error) {

	// Add this user to the subscription system HERE

	response.Success = true
	skylib.LogDebug(*response)
	skylib.Requests.Add(1)
	return
}

func main() {
	flag.Parse()
	sig := &SubscriptionService{}
	skylib.NewAgent().Register(sig).Start().Wait()
}
