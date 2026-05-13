package event

import "encoding/json"

type Event struct {
	Type    string
	Payload json.RawMessage `json:"payload"`
}
