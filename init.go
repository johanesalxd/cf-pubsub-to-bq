package subscribepubsub

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/bigquery"
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
