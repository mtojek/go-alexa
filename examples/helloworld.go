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
