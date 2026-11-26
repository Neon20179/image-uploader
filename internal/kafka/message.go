package kafka

type Message struct {
	Key     []byte
	Value   []byte
	Topic   string
	Headers map[string]string
}
