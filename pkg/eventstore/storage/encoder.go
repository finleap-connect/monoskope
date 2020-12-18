package storage

type Encoder interface {
	Marshal(EventData) ([]byte, error)
	Unmarshal(EventType, []byte) (EventData, error)
	String() string
}
