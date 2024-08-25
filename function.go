package subscribepubsub

import (
	"context"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

// init initializes the Cloud Function and BigQuery client.
// It registers the SubscribePubSub function as a CloudEvent handler
// and initializes the BigQuery client using a sync.Once to ensure
// it's only done once.
func init() {
	// Set up basic logging configuration
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	functions.CloudEvent("SubscribePubSub", SubscribePubSub)
	initOnce.Do(initializeBigQuery)
}

// MessagePublishedData represents the structure of the published message data.
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage represents the structure of a Pub/Sub message.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// SubscribePubSub is the main Cloud Function that handles incoming Pub/Sub messages.
// It processes the CloudEvent, extracts the Pub/Sub message, and inserts it into BigQuery.
// ctx: The context for the function invocation.
// e: The CloudEvent containing the Pub/Sub message.
// Returns an error if any step in the process fails.
func SubscribePubSub(ctx context.Context, e event.Event) error {
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
