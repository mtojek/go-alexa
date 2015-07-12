package alexa

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/context"
)

var (
	ErrRequestExpired = errors.New("request too old to continue (>150s)")
	ErrInvalidAppID   = errors.New("request does not match app ID")
)

func GetEchoRequest(r *http.Request) *EchoRequest {
	return context.Get(r, "echoRequest").(*EchoRequest)
}

type Skill struct {
	AppID string
}

// New creates a new Skill with the given app ID.
func New(appID string) *Skill {
	return &Skill{
		AppID: appID,
	}
}

// verifyJSON decodes the request, verifies it, and stores it in the context of
// the request.
func (s *Skill) HandlerFuncWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var err error

	// Validate the request
	err = ValidateAmazonRequest(r)
	if err != nil {
		httpError(w, err.Error(), "Not Authorized", http.StatusUnauthorized)
		return
	}

	// Parse the request
	echoReq := new(EchoRequest)
	err = json.NewDecoder(r.Body).Decode(&echoReq)
	if err != nil {
		httpError(w, err.Error(), "Bad Request", http.StatusInternalServerError)
		return
	}

	// Check the timestamp
	if !echoReq.VerifyTimestamp(150) && r.URL.Query().Get("_dev") == "" {
		httpError(w, ErrRequestExpired.Error(), "Bad Request", http.StatusInternalServerError)
		return
	}

	// Check the app ID
	if !echoReq.VerifyAppID(s.AppID) {
		httpError(w, ErrInvalidAppID.Error(), "Bad Request", http.StatusInternalServerError)
		return
	}

	// Store the request in the context
	context.Set(r, "echoRequest", echoReq)

	next(w, r)
}

// httpError will respond with the given error and log the message
func httpError(w http.ResponseWriter, logMsg string, err string, errCode int) {
	if logMsg != "" {
		log.Println(logMsg)
	}
	http.Error(w, err, errCode)
}
