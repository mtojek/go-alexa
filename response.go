package alexa

import (
	"encoding/json"
)

type EchoResponse struct {
	Version           string                 `ppson:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          EchoRespBody           `json:"response"`
}

type EchoRespBody struct {
	OutputSpeech     *EchoOutputSpeech `json:"outputSpeech,omitempty"`
	Card             *EchoCard         `json:"card,omitempty"`
	Reprompt         *EchoReprompt     `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	ShouldEndSession bool              `json:"shouldEndSession"`
}

type EchoReprompt struct {
	OutputSpeech EchoOutputSpeech `json:"outputSpeech,omitempty"`
}

type EchoOutputSpeech struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type EchoCard struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

func NewResponse() *EchoResponse {
	er := &EchoResponse{
		Version: "1.0",
		Response: EchoRespBody{
			ShouldEndSession: true,
		},
	}

	return er
}

func (r *EchoResponse) OutputSpeech(text string) *EchoResponse {
	r.Response.OutputSpeech = &EchoOutputSpeech{
		Type: "PlainText",
		Text: text,
	}

	return r
}

func (r *EchoResponse) Card(title string, content string) *EchoResponse {
	r.Response.Card = &EchoCard{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return r
}

func (r *EchoResponse) Reprompt(text string) *EchoResponse {
	r.Response.Reprompt = &EchoReprompt{
		OutputSpeech: EchoOutputSpeech{
			Type: "PlainText",
			Text: text,
		},
	}

	return r
}

func (r *EchoResponse) EndSession(flag bool) *EchoResponse {
	r.Response.ShouldEndSession = flag

	return r
}

func (r *EchoResponse) ToJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}
