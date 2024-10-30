package PubSubToBQ

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/bigquery"
)

var (
	bqInserter *bigquery.Inserter
	initOnce   sync.Once
	initError  error
)

// validateMessage validates the Pub/Sub message data.
// It checks if the message data is empty and returns an error if true.
// msg: The MessagePublishedData containing the Pub/Sub message.
// Returns an error if the message is invalid.
func validateMessage(msg MessagePublishedData) error {
	if len(msg.Message.Data) == 0 {
		return fmt.Errorf("empty message data")
	}
	// Add more validation as needed
	return nil
}

// processMessage unmarshals a message and inserts it into BigQuery
// It returns an error if unmarshaling or insertion fails
func processMessage(ctx context.Context, messageData []byte, inserter *bigquery.Inserter) error {
	var data BigQueryRow
	if err := json.Unmarshal(messageData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if err := inserter.Put(ctx, data); err != nil {
		return fmt.Errorf("failed to insert data into BigQuery: %w", err)
	}

	log.Printf("Successfully inserted data into BigQuery")
	return nil
}
