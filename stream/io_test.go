package stream_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/customerio/homework/stream"
)

func TestProcess(t *testing.T) {
	var expected = []*stream.Record{
		{
			ID:     "111-222-333",
			Type:   "attributes",
			UserID: "user-1",
			Data: map[string]string{
				"name": "john doe",
				"city": "toronto",
			},
			Timestamp: 1234567890,
		},
		{
			ID:     "123-456-789",
			Type:   "event",
			Name:   "signup",
			UserID: "user-2",
			Data: map[string]string{
				"random": "test",
			},
			Timestamp: 1234567890,
		},
	}

	var input bytes.Buffer
	var offset int64
	for i, rec := range expected {
		b, err := json.Marshal(rec)
		if err != nil {
			t.Errorf("error encoding record %d: %v", i, err)
		}
		b = append(b, '\n')
		n, err := input.Write(b)
		if err != nil {
			t.Errorf("error writing record %d: %v", i, err)
		}
		offset += int64(n)
		rec.Position = offset
	}

	ch, err := stream.Process(context.Background(), bytes.NewReader(input.Bytes()))
	if err != nil {
		t.Errorf("error processing data: %v", err)
	}

	var i int
	for rec := range ch {
		if !match(rec, expected[i]) {
			t.Errorf("record %d does not match:\nwant: %#v\nhave: %#v", i, expected[i], rec)
		}
		i++
	}
}

func match(r1, r2 *stream.Record) bool {
	if r1.ID != r2.ID {
		return false
	}
	if r1.Type != r2.Type {
		return false
	}
	if r1.Name != r2.Name {
		return false
	}
	if r1.UserID != r2.UserID {
		return false
	}

	if len(r1.Data) != len(r2.Data) {
		return false
	}
	for k, v1 := range r1.Data {
		v2, found := r2.Data[k]
		if !found {
			return false
		}
		if v1 != v2 {
			return false
		}
	}

	if r1.Timestamp != r2.Timestamp {
		return false
	}
	if r1.Position != r2.Position {
		return false
	}

	return true
}
