package alexa

import (
	"errors"
	"time"
)

type EchoRequest struct {
	Version string      `json:"version"`
	Session EchoSession `json:"session"`
	Request EchoReqBody `json:"request"`
}

type EchoSession struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes struct {
		String map[string]interface{} `json:"string"`
	} `json:"attributes"`
	User struct {
		UserID string `json:"string"`
	} `json:"user"`
}

type EchoReqBody struct {
	Type      string     `json:"type"`
	RequestID string     `json:"requestId"`
	Timestamp string     `json:"timestamp"`
	Intent    EchoIntent `json:"intent,omitempty"`
	Reason    string     `json:"reason,omitempty"`
}

type EchoIntent struct {
	Name  string              `json:"name"`
	Slots map[string]EchoSlot `json:"slots"`
}

type EchoSlot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// VerifyTimestamp returns true if the request is occurred happened within the
// given number of seconds.
func (r *EchoRequest) VerifyTimestamp(seconds int) bool {
	reqTimestamp, err := time.Parse("2006-01-02T15:04:05Z", r.Request.Timestamp)
	if err != nil {
		return false
	}
	return time.Since(reqTimestamp) < time.Duration(seconds)*time.Second
}

// VerifyAppID returns true if the request matches the given app ID.
func (r *EchoRequest) VerifyAppID(myAppID string) bool {
	return r.Session.Application.ApplicationID == myAppID
}

func (r *EchoRequest) GetSessionID() string {
	return r.Session.SessionID
}

func (r *EchoRequest) GetUserID() string {
	return r.Session.User.UserID
}

func (r *EchoRequest) GetRequestType() string {
	return r.Request.Type
}

func (r *EchoRequest) GetIntentName() string {
	if r.GetRequestType() == "IntentRequest" {
		return r.Request.Intent.Name
	}
	return r.GetRequestType()
}

func (r *EchoRequest) GetSlotValue(slotName string) (string, error) {
	if _, ok := r.Request.Intent.Slots[slotName]; ok {
		return r.Request.Intent.Slots[slotName].Value, nil
	}

	return "", errors.New("Slot name not found.")
}

func (r *EchoRequest) AllSlots() map[string]EchoSlot {
	return r.Request.Intent.Slots
}
