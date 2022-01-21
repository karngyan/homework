package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/customerio/homework/datastore"
	"github.com/customerio/homework/serve"
	"github.com/customerio/homework/stream"
	"github.com/customerio/homework/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	// process stream and fetch summarized data
	szData := processStream(ctx, "data/messages.2.data")

	// create datastore
	var ds serve.Datastore
	var err error

	if ds, err = datastore.CreateDatastore(szData.attributes, szData.events); err != nil {
		log.Fatal("failed to create data store, err: ", err)
	}

	// TODO: Can clear the summarized data to free up memory

	// start the server
	if err := serve.ListenAndServe(":1323", ds); err != nil {
		log.Fatal(err)
	}
}

type summarizedData struct {
	// attributes - map of user_id -> most recent Record with summarized attributes
	attributes map[string]stream.Record
	// events - user_id -> event_name -> count
	events map[string]map[string]int
}

func processStream(ctx context.Context, filepath string) summarizedData {
	var totRecordsProcessed = 0
	var start = time.Now()
	var file *os.File
	var err error

	sd := summarizedData{
		attributes: make(map[string]stream.Record),
		events:     make(map[string]map[string]int),
	}

	if file, err = os.Open(filepath); err != nil {
		log.Fatal(fmt.Sprintf("failed to open file, error: %v", err))
	}

	if ch, err := stream.Process(ctx, file); err == nil {

		// dupEvents hash to keep track of duplicate events - event_id -> bool
		var dupEvents = make(map[string]bool)

		for rec := range ch {
			totRecordsProcessed++
			switch rec.Type {
			case "event":
				if _, prs := dupEvents[rec.ID]; prs {
					continue
				}

				if _, prs := sd.events[rec.UserID]; !prs {
					sd.events[rec.UserID] = make(map[string]int)
				}

				dupEvents[rec.ID] = true
				sd.events[rec.UserID][rec.Name] += 1

			case "attributes":
				// attributes are merged to prevent last-write-wins scenario
				if xrecord, prs := sd.attributes[rec.UserID]; prs {
					if rec.Timestamp >= xrecord.Timestamp {
						// recent record so overwrite any common keys in `xrecord.Data`
						rec.Data = utils.MergeMaps(rec.Data, xrecord.Data, true)

					} else {
						// rec is old record, so keep `xrecord.Data & xrecord.Timestamp` but add any new keys in rec.Data
						rec.Data = utils.MergeMaps(xrecord.Data, rec.Data, true)
						rec.Timestamp = xrecord.Timestamp
						rec.ID = xrecord.ID
						rec.Position = xrecord.Position
					}
				}

				sd.attributes[rec.UserID] = *rec
			}
		}
		if err := ctx.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("stream processing failed", err)
	}

	duration := time.Now().Sub(start)
	log.Println("time taken to process: ", duration)
	log.Println("total records processed: ", totRecordsProcessed)

	return sd
}
