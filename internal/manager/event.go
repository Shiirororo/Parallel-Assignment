package manager

import "encoding/json"

type Event struct {
	Type    string
	Payload json.RawMessage `json:"payload"`
}
