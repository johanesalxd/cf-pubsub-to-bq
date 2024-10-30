PubSub to BigQuery using Cloud Run Functions
-----------------------------
Details TBA

# How to run
## Run locally
TBA

## Test locally
TBA

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
