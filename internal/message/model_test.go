package message

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestFromJSON(t *testing.T) {
	tests := []struct {
		explanation     string
		input           string
		messageExpected *Message
		errExpected     error
	}{
		{
			"valid message",
			`{"from":"gopher","msg":"hello","time":"2025-04-20T12:00:00Z"}`,
			&Message{
				From: "gopher",
				Msg:  "hello",
				Time: time.Date(2025, 4, 20, 12, 0, 0, 0, time.UTC),
			},
			nil,
		},
		{
			"invalid JSON",
			`{"from":"gopher",`,
			nil,
			ErrMalformedMessage,
		},
		{
			"unicode and emoji",
			`{"from":"bot","msg":"👋🏼 привет","time":"2025-04-20T12:00:00Z"}`,
			&Message{
				From: "bot",
				Msg:  "👋🏼 привет",
				Time: time.Date(2025, 4, 20, 12, 0, 0, 0, time.UTC),
			},
			nil,
		},
		{
			"missing fields",
			`{"from":"ghost"}`,
			nil,
			ErrMissingField,
		},
		{
			"empty JSON",
			`{}`,
			nil,
			ErrMissingField,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s [%s]", tt.explanation, tt.input), func(t *testing.T) {
			message, err := FromJSONString(tt.input)

			if got, want := err, tt.errExpected; !errors.Is(got, want) {
				t.Fatalf("err=%T (%v), want=%T (%v)", got, got, want, want)
			}

			if got, want := message, tt.messageExpected; !got.Equal(want) {
				t.Errorf("message=%v, want=%v", got, want)
			}
		})
	}
}

func TestToBytes(t *testing.T) {
	tests := []struct {
		explanation  string
		input        *Message
		dataExpected []byte
		errExpected  error
	}{
		{
			"valid message",
			&Message{
				From: "gopher",
				Msg:  "hello",
				Time: time.Date(2025, 4, 20, 12, 0, 0, 0, time.UTC),
			},
			[]byte(`{"from":"gopher","msg":"hello","time":"2025-04-20T12:00:00Z"}`),
			nil,
		},
		{
			"empty message",
			&Message{
				From: "",
				Msg:  "",
				Time: time.Time{},
			},
			nil,
			ErrMissingField,
		},
		{
			"message with special characters",
			&Message{
				From: "bot",
				Msg:  "👋🏼 привет",
				Time: time.Date(2025, 4, 20, 12, 0, 0, 0, time.UTC),
			},
			[]byte(`{"from":"bot","msg":"👋🏼 привет","time":"2025-04-20T12:00:00Z"}`),
			nil,
		},
		{
			"message with only from field",
			&Message{
				From: "ghost",
				Msg:  "",
				Time: time.Time{},
			},
			nil,
			ErrMissingField,
		},
		{
			"missing Msg field",
			&Message{
				From: "ghost",
				Msg:  "",
				Time: time.Time{},
			},
			nil,
			ErrMissingField,
		},
		{
			"missing Time field",
			&Message{
				From: "ghost",
				Msg:  "hello",
				Time: time.Time{},
			},
			nil,
			ErrMissingField,
		},
	}

	for _, tt := range tests {
		t.Run(tt.explanation, func(t *testing.T) {
			data, err := tt.input.ToBytes()
			if got, want := err, tt.errExpected; !errors.Is(got, want) {
				t.Fatalf("err=%T (%v), want=%T (%v)", got, got, want, want)
			}

			if got, want := data, tt.dataExpected; !bytes.Equal(got, want) {
				t.Errorf("data=%v, want=%v", string(got), string(want))
			}
		})
	}
}
