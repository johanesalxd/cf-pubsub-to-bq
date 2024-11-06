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
  * Update `TOPIC_NAME`, `YOUR_PROJECT_ID`, `REGION_NAME`, `SCHEMA_LOCATION` and `QPS` as per requirements

### Pub/Sub
```
gcloud pubsub topics create TOPIC_NAME
```
```
gcloud pubsub topics create PubSubToBQ
```

### Dataflow
```
gcloud storage cp data-generator/source.json gs://mybucket
```
```
gcloud dataflow flex-template run data-generator-pubsub-to-bq \
    --project=YOUR_PROJECT_ID \
    --region=REGION_NAME \
    --template-file-gcs-location=gs://dataflow-templates-REGION_NAME/latest/flex/Streaming_Data_Generator \
    --parameters \
schemaLocation=SCHEMA_LOCATION,\
qps=QPS,\
topic=projects/YOUR_PROJECT_ID/topics/PubSubToBQ
```
```
gcloud dataflow flex-template run data-generator-pubsub-to-bq \
    --project=YOUR_PROJECT_ID \
    --region=us-central1 \
    --template-file-gcs-location=gs://dataflow-templates-us-central1/latest/flex/Streaming_Data_Generator \
    --parameters \
schemaLocation=gs://mybucket/source.json,\
qps=1,\
topic=projects/YOUR_PROJECT_ID/topics/PubSubToBQ
```

## BigQuery model and example
* Update `YOUR_DATASET_NAME` and `YOUR_TABLE_NAME` accordingly
* Table structure is aligned with `model.go`
```
bq query --nouse_legacy_sql \
'CREATE OR REPLACE TABLE
  YOUR_DATASET_NAME.YOUR_TABLE_NAME ( event_type string,
    timestamp timestamp,
    player_id string,
    game_version string,
    device_id string,
    location string )
PARTITION BY
  DATE(timestamp);'
```
```
bq query --nouse_legacy_sql \
'SELECT
  DATE(timestamp) AS pt,
  event_type,
  COUNT(1) AS cnt
FROM
  YOUR_DATASET_NAME.YOUR_TABLE_NAME
GROUP BY
  1,
  2;'
```

## Run on Cloud Function
Notes: update `trigger-topic` and `.env.yaml` accordingly
```
gcloud functions deploy cf-pubsub-to-bq \
    --gen2 \
    --runtime=go122 \
    --region=us-central1 \
    --source=. \
    --entry-point=PubSubToBQ \
    --trigger-topic=PubSubToBQ \
    --allow-unauthenticated \
    --env-vars-file=.env.yaml
```

# Additional notes
TBA

## Related links
TBA