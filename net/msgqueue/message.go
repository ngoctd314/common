package queue

type Message struct {
	kafakaMessage *KafkaMessage
}

type KafkaMessage struct {
	Key   []byte
	Value []byte
}

func NewKafkaMessage(msg *KafkaMessage) *Message {
	return &Message{
		kafakaMessage: msg,
	}
}
