package stream

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Record struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	UserID    string            `json:"user_id"`
	Data      map[string]string `json:"data"`
	Timestamp int64             `json:"timestamp"`

	// Position in the input stream where this record lives.
	Position int64 `json:"-"`
}

// Process returns a channel to which a stream of records are sent. Reading starts at
// the current seek offset in the file. The channel is closed when no more records are available.
// If the context completes, reading is prematurely terminated.
func Process(ctx context.Context, f io.ReadSeeker) (<-chan *Record, error) {

	if f == nil {
		return nil, fmt.Errorf("must supply a valid io.ReadSeeker (probably an *os.File) to the stream.Process function")
	}

	offset, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	ch := make(chan *Record)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(f)
		scanner.Split(func(data []byte, atEof bool) (advance int, token []byte, err error) {
			advance, token, err = bufio.ScanLines(data, atEof)
			if err == nil && token != nil {
				offset += int64(advance)
			}
			return advance, token, err
		})
		for scanner.Scan() {
			rec := &Record{
				Position: offset,
			}
			if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
				log.Println("json.Unmarshal failed", err)
				continue
			}
			select {
			case _ = <-ctx.Done():
				return
			case ch <- rec:
			}
		}
	}()
	return ch, nil
}
