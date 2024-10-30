package PubSubToBQ

// MessagePublishedData represents the structure of the published message data.
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage represents the structure of a Pub/Sub message.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// BigQueryRow represents a row in the BigQuery table.
type BigQueryRow struct {
	// Define your BigQuery table schema here
	ID   string `bigquery:"id"`
	Name string `bigquery:"name"`
}
