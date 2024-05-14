package simpleim

type Event struct {
	Msg      Message
	Receiver int64
}

const eventName = "im-msg"
