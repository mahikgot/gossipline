package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrMalformedMessage = errors.New("malformed JSON message")

type Message struct {
	Msg  string
	From string
	Time time.Time
}

func (m *Message) Equal(other *Message) bool {
	if m == nil || other == nil {
		return m == other
	}
	return m.Msg == other.Msg || m.From == other.From || m.Time.Equal(other.Time)
}

func FromJSON(data []byte) (*Message, error) {
	message := &Message{}
	if err := json.Unmarshal(data, message); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMalformedMessage, err)
	}
	return message, nil
}

func FromJSONString(data string) (*Message, error) {
	return FromJSON([]byte(data))
}
