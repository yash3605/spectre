package models

type Tab int
type State string

const (
	OSINT Tab = iota
	Infosys
	Entity
	TabCount
)

const (
	StateIdle    State = "idle"
	StateLoading State = "loading"
	StateError   State = "error"
	StateSuccess State = "success"
)

type Result struct {
	Title  string
	Data   map[string]string
	Status State
}
