package types

const Evlogmsg = "evlogmsg"
const EvWorkerStopped = "EvWorkerStopped"

type LogMsg struct {
	Msg   string
	Level string
}
