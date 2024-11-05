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
	EventType   string `bigquery:"event_type" json:"event_type"`
	Timestamp   string `bigquery:"timestamp" json:"timestamp"`
	PlayerID    string `bigquery:"player_id" json:"player_id"`
	GameVersion string `bigquery:"game_version" json:"game_version"`
	DeviceID    string `bigquery:"device_id" json:"device_id"`
	Location    string `bigquery:"location" json:"location"`
}
