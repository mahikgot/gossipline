package net

import (
	"errors"
	"net"

	"github.com/mahikgot/gossipline/internal/message"
)

var ErrWriteFail = errors.New("failed to send mesage")

type Transmitter struct {
	conn net.Conn
}

func (t *Transmitter) Send(m *message.Message) error {
	data, err := m.ToBytes()
	if err != nil {
		return err
	}
	_, err = t.conn.Write(data)

	return err
}

func (t *Transmitter) Recieve(ch chan []byte) error {
	for {
		data := make([]byte, 512)
		n, err := t.conn.Read(data)
		if err != nil {
			return err
		}
		ch <- data[:n]
	}
}

func RecieveMessages(ch chan []byte) {
	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	transmitter := Transmitter{conn}
	go transmitter.Recieve(ch)
}
