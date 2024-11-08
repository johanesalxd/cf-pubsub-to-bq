PubSub to BigQuery using Cloud Run Functions
-----------------------------
Details TBA

# How to test
## Input model and example
* You can use the included *data generator* or build yourself from [here](https://github.com/vincentrussell/json-data-generator/tree/json-data-generator-1.16)
```
java -jar data-generator/json-data-generator-1.16-standalone.jar -s data-generator/source.json
{
    "event_type": "game_match",
    "timestamp": "2024-11-05T11:12:50Z",
    "player_id": "cvj50okbHO",
    "game_version": "2.5",
    "device_id": "android",
    "location": "Portugal"
}
```

## Test locally
TBA

# How to run
## Data Generator setup
* Using Dataflow data generator from [here](https://cloud.google.com/dataflow/docs/guides/templates/provided/streaming-data-generator)
  * Update `TOPIC_NAME`, `PROJECT_ID`, `REGION_NAME`, `SCHEMA_BUCKET_LOCATION` and `QPS` as per requirements

### Create Pub/Sub Topic and Subscription
```
gcloud pubsub topics create TOPIC_NAME
```
```
gcloud pubsub subscriptions create TOPIC_NAME-sub --topic=TOPIC_NAME
```

### Data Generator deployment
* It may take couple of minutes to complete. You can check in the [Dataflow console](https://console.cloud.google.com/dataflow/) for the latest status 
```
gcloud storage cp data-generator/source.json gs://SCHEMA_BUCKET_LOCATION
```
```
gcloud dataflow flex-template run data-generator-pubsub-to-bq \
    --project=PROJECT_ID \
    --region=REGION_NAME \
    --template-file-gcs-location=gs://dataflow-templates-REGION_NAME/latest/flex/Streaming_Data_Generator \
    --parameters \
schemaLocation=gs://SCHEMA_BUCKET_LOCATION,\
qps=QPS,\
topic=projects/PROJECT_ID/topics/TOPIC_NAME
```
For example:
```
gcloud dataflow flex-template run data-generator-pubsub-to-bq \
    --project=PROJECT_ID \
    --region=us-central1 \
    --template-file-gcs-location=gs://dataflow-templates-us-central1/latest/flex/Streaming_Data_Generator \
    --parameters \
schemaLocation=gs://mybucket/source.json,\
qps=1,\
topic=projects/PROJECT_ID/topics/TOPIC_NAME
```

### Checkpoint
* Update `TOPIC_NAME` as per requirements
* Check if the generated messages from the previous steps has been delivered successfully to Pub/Sub Topic
```
gcloud pubsub subscriptions pull TOPIC_NAME-sub --no-auto-ack --format=json
```

## Data ingestion setup
### Create BigQuery table
* Update `DATASET_NAME` and `TABLE_NAME` accordingly
* Table structure is aligned with `model.go`
```
bq query --nouse_legacy_sql \
'CREATE OR REPLACE TABLE
  DATASET_NAME.TABLE_NAME ( event_type string,
    timestamp timestamp,
    player_id string,
    game_version string,
    device_id string,
    location string )
PARTITION BY
  DATE(timestamp);'
```

### Cloud Run functions deployment
* Update `TOPIC_NAME` and `.env.yaml` accordingly 
```
gcloud functions deploy cf-pubsub-to-bq \
    --gen2 \
    --runtime=go122 \
    --region=us-central1 \
    --source=. \
    --entry-point=PubSubToBQ \
    --trigger-topic=TOPIC_NAME \
    --allow-unauthenticated \
    --env-vars-file=.env.yaml
```

### Checkpoint
* Update `DATASET_NAME` and `TABLE_NAME` accordingly
```
bq query --nouse_legacy_sql \
'SELECT
  DATE(timestamp) AS pt,
  event_type,
  COUNT(1) AS cnt
FROM
  DATASET_NAME.TABLE_NAME
GROUP BY
  1,
  2;'
```

# Additional notes
TBA

## Related links
TBA