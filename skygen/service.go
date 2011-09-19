//Copyright (c) 2011 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

const serviceTemplate = `package main

import (
	"github.com/bketelsen/skynet/skylib"
	"os"
	"flag"
	"<%PackageName%>"
)

type <%ServiceName%>Service struct {

}

func (*<%ServiceName%>Service) Process<%ServiceName%>(m <%PackageName%>.<%ServiceName%>Request, response *<%PackageName%>.<%ServiceName%>Response) (err os.Error) {

	//Process the message here
	skylib.LogDebug(m.EmailAddress)
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
`
