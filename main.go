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

	var ds serve.Datastore
	var file *os.File
	var err error

	if file, err = os.Open("data/messages.2.data"); err != nil {
		log.Fatal(fmt.Sprintf("failed to open file, error: %v", err))
	}

	processStart := time.Now()
	// summarizedAttributes - map of user_id -> most recent Record with summarized attributes
	var summarizedAttributes = make(map[string]stream.Record)
	// summarizedEvents - user_id -> event_name -> count
	var summarizedEvents = make(map[string]map[string]int)

	var totalRecordsProcessed = 0

	if ch, err := stream.Process(ctx, file); err == nil {

		// dupEvents hash to keep track of duplicate events - event_id -> bool
		var dupEvents = make(map[string]bool)

		for rec := range ch {
			totalRecordsProcessed++
			switch rec.Type {
			case "event":
				if _, prs := dupEvents[rec.ID]; prs {
					continue
				}

				if _, prs := summarizedEvents[rec.UserID]; !prs {
					summarizedEvents[rec.UserID] = make(map[string]int)
				}

				dupEvents[rec.ID] = true
				summarizedEvents[rec.UserID][rec.Name] += 1

			case "attributes":
				// attributes are merged to prevent last-write-wins scenario
				if xrecord, prs := summarizedAttributes[rec.UserID]; prs {
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

				summarizedAttributes[rec.UserID] = *rec
			}
		}
		if err := ctx.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("stream processing failed", err)
	}

	processDuration := time.Now().Sub(processStart)
	log.Println("time taken to process: ", processDuration)
	log.Println("total records processed: ", totalRecordsProcessed)

	// create datastore
	if ds, err = datastore.CreateDatastore(summarizedAttributes, summarizedEvents); err != nil {
		log.Fatal("failed to create data store, err: ", err)
	}

	// start the server
	if err := serve.ListenAndServe(":1323", ds); err != nil {
		log.Fatal(err)
	}
}
