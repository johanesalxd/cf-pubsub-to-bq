PubSub to BigQuery using Cloud Run Functions
-----------------------------
Details TBA

# How to run
## Run locally
TBA

# How to test
## Input model and example
* You can use the included *data generator* or build yourself from [here](https://github.com/vincentrussell/json-data-generator/tree/json-data-generator-1.16)
```
# source.json
{
    "event_type": "{{random("login","logout","game_search","game_match","game_join","game_disconnect","game_reconnect")}}",
    "timestamp": "{{date("yyyy-MM-dd'T'HH:mm:ss'Z'")}}",
    "player_id": "{{alphaNumeric(10)}}",
    "game_version": "{{random("2.0","2.5","3.0")}}",
    "device_id": "{{random("android","ios")}}",
    "location": "{{country()}}"
}
```
```
# java -jar data-generator/json-data-generator-1.16-standalone.jar -s source.json
{
    "event_type": "9ddde82e-c578-4ac2-9ddb-17f063659d88",
    "timestamp": "2024-11-05T00:26:26Z",
    "player_id": "uKb33JHT7m",
    "game_version": "2.0",
    "device_id": "ios",
    "location": "Argentina"
}
```
## BigQuery model and example
Notes: update `YOUR_DATASET_NAME` and `YOUR_TABLE_NAME` accordingly
```
CREATE OR REPLACE TABLE
  YOUR_DATASET_NAME.YOUR_TABLE_NAME ( event_type string,
    timestamp timestamp,
    player_id string,
    game_version string,
    device_id string,
    location string )
PARTITION BY
  DATE(timestamp);
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

## Run on Cloud Run
TBA

# Additional notes
TBA

## Related links
* https://cloud.google.com/functions/docs/tutorials/pubsub#objectives
* https://cloud.google.com/functions/docs/local-development
