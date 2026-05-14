package event

import "encoding/json"

type Event struct {
	Type    string
	Request Request
}

type Request struct {
	OriginID string
	Payload  json.RawMessage `json:"payload"`
}
