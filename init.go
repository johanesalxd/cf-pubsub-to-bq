package PubSubToBQ

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type Config struct {
	ProjectID string
	DatasetID string
	TableID   string
}

// defaultConfig creates and returns a Config struct with default values
// from environment variables.
func defaultConfig() *Config {
	return &Config{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		DatasetID: os.Getenv("DATASET_ID"),
		TableID:   os.Getenv("TABLE_ID"),
	}
}

// loadConfig loads the application configuration from environment variables.
// It returns an error if any required environment variable is not set.
func loadConfig() (*Config, error) {
	config := defaultConfig()

	if config.ProjectID == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT environment variable is not set")
	}

	if config.DatasetID == "" {
		return nil, fmt.Errorf("DATASET_ID environment variable is not set")
	}

	if config.TableID == "" {
		return nil, fmt.Errorf("TABLE_ID environment variable is not set")
	}

	return config, nil
}

// init initializes the Cloud Function and BigQuery client.
// It registers the PubSubToBQ function as a CloudEvent handler
// and initializes the BigQuery client using a sync.Once to ensure
// it's only done once.
func init() {
	// Set up basic logging configuration
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	functions.CloudEvent("PubSubToBQ", PubSubToBQ)
	initOnce.Do(initializeBigQuery)
}

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

// initializeBigQueryClient creates and initializes a BigQuery client.
// It takes a context and a Config struct, and returns a BigQuery client and an error.
func initializeBigQueryClient(ctx context.Context, config *Config) (*bigquery.Client, error) {
	// Initialize BigQuery client
	bqClient, err := bigquery.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create BigQuery client: %v", err)
	}

	return bqClient, nil
}
