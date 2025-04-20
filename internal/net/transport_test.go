package net

import (
	"bytes"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/mahikgot/gossipline/internal/message"
)

func TestTransport(t *testing.T) {
	tests := []struct {
		explanation  string
		input        message.Message
		dataExpected []byte
		errExpected  error
	}{
		{
			"valid message",
			message.Message{From: "gopher", Msg: "hello", Time: time.Date(2025, 4, 20, 12, 0, 0, 0, time.UTC)},
			[]byte(`{"from":"gopher","msg":"hello","time":"2025-04-20T12:00:00Z"}`),
			nil,
		},
		{
			"empty message",
			message.Message{From: "", Msg: ""},
			nil,
			message.ErrMissingField,
		},
		{
			"message with special characters",
			message.Message{From: "bot", Msg: "👋🏼 привет", Time: time.Date(2025, 4, 20, 12, 0, 0, 0, time.UTC)},
			[]byte(`{"from":"bot","msg":"👋🏼 привет","time":"2025-04-20T12:00:00Z"}`),
			nil,
		},
		{
			"error in writing",
			message.Message{From: "errorTest", Msg: "fail"},
			nil,
			message.ErrMissingField,
		},
	}

	for _, tt := range tests {
		t.Run(tt.explanation, func(t *testing.T) {
			mockWriterTo := &mockWriterTo{}

			transmitter := &Transmitter{
				conn: mockWriterTo,
				dst:  nil,
			}

			err := transmitter.Send(&tt.input)

			if got, want := err, tt.errExpected; !errors.Is(got, want) {
				t.Fatalf("err=%v, want=%v", got, want)
			}

			if got := mockWriterTo.Bytes; !bytes.Equal(got, tt.dataExpected) {
				t.Errorf("data=%s, want=%s", string(got), string(tt.dataExpected))
			}
		})
	}
}

type mockWriterTo struct {
	Bytes []byte
}

func (m *mockWriterTo) WriteTo(p []byte, addr net.Addr) (int, error) {
	m.Bytes = p
	return len(p), nil
}
