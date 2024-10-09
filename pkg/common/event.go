package common

const (
	ACTION_ADD    = "ADD"
	ACTION_REMOVE = "REMOVE"
)

type Event struct {
	Action string `json:"action"`
	File   string `json:"file"`
}