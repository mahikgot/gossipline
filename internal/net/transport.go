package net

import (
	"errors"
	"net"

	"github.com/mahikgot/gossipline/internal/message"
)

var ErrWriteFail = errors.New("failed to send mesage")

type Transmitter struct {
	conn writerTo
	dst  net.Addr
}

func (t *Transmitter) Send(m *message.Message) error {
	data, err := m.ToBytes()
	if err != nil {
		return err
	}
	_, err = t.conn.WriteTo(data, t.dst)

	return err
}

type Sender interface {
	Send(m *message.Message) error
}

type Reciever interface {
	Recieve(data []byte) error
}

type writerTo interface {
	WriteTo([]byte, net.Addr) (int, error)
}
