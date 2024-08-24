package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
)

type Config struct {
	ProjectID      string
	DatasetID      string
	TableID        string
	SubscriptionID string
	NumWorkers     int
}

func defaultConfig() *Config {
	return &Config{
		ProjectID:      os.Getenv("GOOGLE_CLOUD_PROJECT"),
		DatasetID:      os.Getenv("DATASET_ID"),
		TableID:        os.Getenv("TABLE_ID"),
		SubscriptionID: os.Getenv("SUBSCRIPTION_ID"),
		NumWorkers:     10, // Default number of workers
	}
}

// loadConfig loads the application configuration from environment variables
// and applies any provided functional options
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

	if config.SubscriptionID == "" {
		return nil, fmt.Errorf("SUBSCRIPTION_ID environment variable is not set")
	}

	return config, nil
}

// initializeClients creates and initializes Pub/Sub and BigQuery clients
// It returns both clients or an error if initialization fails
func initializeClients(ctx context.Context, config *Config) (*pubsub.Client, *bigquery.Client, error) {
	// Initialize Pub/Sub client
	pubsubClient, err := pubsub.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create pubsub client: %v", err)
	}

	// Initialize BigQuery client
	bqClient, err := bigquery.NewClient(ctx, config.ProjectID)
	if err != nil {
		pubsubClient.Close()
		return nil, nil, fmt.Errorf("failed to create BigQuery client: %v", err)
	}

	return pubsubClient, bqClient, nil
}
