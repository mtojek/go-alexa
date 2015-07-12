## `go-alexa`

A simple Go framework to quickly create an Amazon Alexa Skills web service.

### What?

After beta testing the Amazon Echo (and it's voice assistant Alexa) for several months, Amazon has released the product to the public and created an SDK for developers to add new "Alexa Skills" to the product.

You can see the SDK documentation here: [developer.amazon.com/public/solutions/alexa/alexa-skills-kit](https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit) but in short, a developer can make a web service that allows a user to say: _Alexa, ask [your service] to [some action your service provides]_

### Requirements

Amazon has a list of requirements to get a new Skill up and running

1. Creating your new Skill on their Development Dashboard populating it with details and example phrases. That process is documented [here](https://developer.amazon.com/appsandservices/solutions/alexa/alexa-skills-kit/docs/defining-the-voice-interface)
2. A lengthy request validation proces. Documented [here](https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/developing-an-alexa-skill-as-a-web-service#Verifying that the Request was Sent by Alexa)
3. A formatted JSON response.
4. SSL connection required, even for development.

### How `go-alexa` Helps

`go-alexa` takes care of #2 and #3 for you so you can concentrate on #1 and coding your app. (#4 is what it is. See the section on SSL below.)

### An Example App

Creating an Alexa Skill web service is easy with `go-alexa`. Simply import the project as any other Go project, define your app, and write your endpoint. All the security checks, request parsing, and response formating are done for you.

Here's a simple, but complete web service example:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/tbuckley/go-alexa"
)

var (
	myskill *alexa.Skill
)

func main() {
	myskill = alexa.New("AMAZON APP ID") // Create a new skill with your app id

	mux := http.NewServeMux()
	mux.HandleFunc("/", HomePage)
	mux.HandleFunc("/echo/helloworld", EchoHelloWorld)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home Page!")
}

func EchoHelloWorld(w http.ResponseWriter, r *http.Request) {
	// Use HandlerFuncWithNext to wrap existing functions or with a framework
	// like Negroni
	myskill.HandlerFuncWithNext(w, r, func(w http.ResponseWriter, r *http.Request) {
		// Get the Echo request object
		echoReq := alexa.GetEchoRequest(r)

		if echoReq.GetRequestType() == "IntentRequest" || echoReq.GetRequestType() == "LaunchRequest" {
			// Create a response
			echoResp := alexa.NewResponse()
			echoResp.OutputSpeech("Hello world from my new Echo test app!")
			echoResp.Card("Hello World", "This is a test card.")

			json, _ := echoResp.ToJSON()
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			w.Write(json)
		}
	})
}

```

Details:
* Your handler is a regular `net/http` handler so you have full access to the request and ResponseWriter.
* The JSON from the Echo request is already parsed for you. Grab it by calling `alexa.GetEchoRequest(r *http.Request)`.
* You generate the Echo Response by using the EchoResponse struct that has methods to generate each part and, when ready, to return it as a JSON string that you can send back as the response.

### The SSL Requirement

Amazon requires an SSL connection for all steps in the Skill process, even local development (which still gets requests from the Echo web service). Amazon is pushing their AWS Lamda service that takes care of SSL for you, but Go isn't an option on Lamda. What I've done personally is put Nginx in front of my Go app and let Nginx handle the SSL (a self-signed cert for development and a real cert when pushing to production). More information here on  [nginx.com](https://www.nginx.com/blog/nginx-ssl/).

### Contributors

* Mike Flynn ([@thatmikeflynn](https://twitter.com/thatmikeflynn))
* Tom Buckley ([tbuckley](https://github.com/tbuckley))
