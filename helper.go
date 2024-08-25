package subscribepubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"cloud.google.com/go/bigquery"
)

type BigQueryRow struct {
	// Define your BigQuery table schema here
	ID   string `bigquery:"id"`
	Name string `bigquery:"name"`
}

var (
	bqInserter *bigquery.Inserter
	initOnce   sync.Once
	initError  error
)

// initializeBigQuery initializes the BigQuery client and inserter.
// It loads the configuration, creates a BigQuery client, and sets up the inserter.
// This function is called once using sync.Once to ensure single initialization.
func initializeBigQuery() {
	ctx := context.Background()

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		initError = fmt.Errorf("failed to load config: %v", err)
		return
	}

	// Initialize BigQuery client
	bqClient, err := initializeBigQueryClient(ctx, config)
	if err != nil {
		initError = fmt.Errorf("failed to initialize BigQuery client: %v", err)
		return
	}

	// Set up BigQuery inserter
	dataset := bqClient.Dataset(config.DatasetID)
	table := dataset.Table(config.TableID)
	bqInserter = table.Inserter()

	log.Println("BigQuery inserter initialized successfully")
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
