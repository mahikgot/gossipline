package net

import (
	"errors"
	"net"

	"github.com/mahikgot/gossipline/internal/message"
)

var ErrWriteFail = errors.New("failed to send mesage")

type Transmitter struct {
	conn net.Conn
	dst  net.Addr
}

func (t *Transmitter) Send(m *message.Message) error {
	data, err := m.ToBytes()
	if err != nil {
		return err
	}
	_, err = t.conn.Write(data)

	return err
}
