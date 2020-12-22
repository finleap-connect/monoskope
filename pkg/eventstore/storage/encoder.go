package storage

// Encoder is an interface for event data encoders for storing and retrieving
// event data from the event store
type Encoder interface {
	Marshal(EventData) ([]byte, error)
	Unmarshal(EventType, []byte) (EventData, error)
	String() string
}
