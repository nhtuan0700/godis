package core

type Command struct {
	Cmd  string
	Args []string
}

const (
	CMD_PING = "PING"
)
