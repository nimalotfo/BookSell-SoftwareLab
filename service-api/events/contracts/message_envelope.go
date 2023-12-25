package contracts

type MessageEnvelope struct {
	EventName string      `json:"event_name"`
	Payload   interface{} `json:"payload"`
}
