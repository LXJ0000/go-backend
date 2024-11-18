package chat

type Chat interface {
	NormalChat(prompt, request string) (response string, err error)
	StreamChat(prompt, request string) (response <-chan string, err error)
}
