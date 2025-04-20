package message

import (
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
			handle, err := FromJSONString(tt.input)

			if got, want := err, tt.errExpected; !errors.Is(got, want) {
				t.Fatalf("err=%T (%v), want=%T (%v)", got, got, want, want)
			}

			if got, want := handle, tt.messageExpected; !got.Equal(want) {
				t.Errorf("handle=%v, want=%v", got, want)
			}
		})
	}
}
