package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
)

type BigQueryRow struct {
	// Define your BigQuery table schema here
	ID   string `bigquery:"id"`
	Name string `bigquery:"name"`
}

// processMessage unmarshals a Pub/Sub message and inserts it into BigQuery
// It returns an error if unmarshaling or insertion fails
func processMessage(ctx context.Context, msg *pubsub.Message, inserter *bigquery.Inserter) error {
	var data BigQueryRow
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if err := inserter.Put(ctx, data); err != nil {
		return fmt.Errorf("failed to insert data into BigQuery: %w", err)
	}

	log.Printf("Successfully inserted data into BigQuery")
	return nil
}

// worker continuously processes messages from the jobs channel
// It acknowledges successful messages and negative-acknowledges failed ones
func worker(ctx context.Context, jobs <-chan *pubsub.Message, inserter *bigquery.Inserter) error {
	for {
		select {
		case msg := <-jobs:
			if err := processMessage(ctx, msg, inserter); err != nil {
				log.Printf("Failed to process message: %v", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
