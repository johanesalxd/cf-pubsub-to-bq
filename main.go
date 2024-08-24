package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
	"golang.org/x/sync/errgroup"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pubsubClient, bqClient, err := initializeClients(ctx, config)
	if err != nil {
		log.Fatalf("Failed to initialize clients: %v", err)
	}
	defer pubsubClient.Close()
	defer bqClient.Close()

	// Set up Pub/Sub subscription
	sub := pubsubClient.Subscription(config.SubscriptionID)

	// Set up BigQuery inserter
	dataset := bqClient.Dataset(config.DatasetID)
	table := dataset.Table(config.TableID)
	inserter := table.Inserter()

	// Set up signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Use errgroup for concurrent operations
	g, ctx := errgroup.WithContext(ctx)
	log.Printf("Waiting for messages...")

	// Create a channel for distributing work
	jobs := make(chan *pubsub.Message, 100) // Buffer size can be adjusted

	// Start worker pool
	for i := 0; i < config.NumWorkers; i++ {
		g.Go(func() error {
			return worker(ctx, jobs, inserter)
		})
	}

	// Start the Pub/Sub subscription receiver
	g.Go(func() error {
		return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			// Distribute messages to workers or handle shutdown
			select {
			case jobs <- msg:
				// Message sent to a worker
			case <-ctx.Done():
				msg.Nack() // Nack if we're shutting down
				return
			}
		})
	})

	// Wait for shutdown signal or error
	select {
	case <-stop:
		log.Println("Shutting down...")
	case <-ctx.Done():
		log.Println("Shutting down due to error...")
	}

	// Initiate graceful shutdown
	cancel()
	if err := g.Wait(); err != nil && err != context.Canceled {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Shutdown complete")
}
