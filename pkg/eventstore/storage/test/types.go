package test

type TestEventData struct {
	Hello string `json:",omitempty"`
}

type TestEventDataExtened struct {
	Hello string `json:",omitempty"`
	World string `json:",omitempty"`
}
