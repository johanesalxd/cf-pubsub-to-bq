package PubSubToBQ

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudevents/sdk-go/v2/event"
)

// PubSubToBQ is the main Cloud Function that handles incoming Pub/Sub messages.
// It processes the CloudEvent, extracts the Pub/Sub message, and inserts it into BigQuery.
// ctx: The context for the function invocation.
// e: The CloudEvent containing the Pub/Sub message.
// Returns an error if any step in the process fails.
func PubSubToBQ(ctx context.Context, e event.Event) error {
	if initError != nil {
		log.Printf("Initialization error: %v", initError)
		return fmt.Errorf("initialization error: %w", initError)
	}

	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		log.Printf("Error parsing Pub/Sub message: %v", err)
		return fmt.Errorf("error parsing Pub/Sub message: %w", err)
	}

	if err := validateMessage(msg); err != nil {
		log.Printf("Invalid message: %v", err)
		return fmt.Errorf("invalid message: %w", err)
	}

	if err := processMessage(ctx, msg.Message.Data, bqInserter); err != nil {
		log.Printf("Error processing message: %v", err)
		return fmt.Errorf("error processing message: %w", err)
	}

	log.Println("Message processed successfully")
	return nil
}