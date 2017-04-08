#!/bin/bash

# Google
export APP_NAME=$(basename "$PWD")
export PROJECT_ID="mchmarny-dev"
export PUBSUB_TW_TOPIC="tweets"
export PUBSUB_RZ_TOPIC="results"
export PUBSUB_TW_SUB="${PUBSUB_TW_TOPIC}-events"
export SERVICE_ACCOUNT_NAME="${APP_NAME}-user"


echo "Canceling Dataflow jobs..."
for id in  $(gcloud beta dataflow jobs list --format='value(JOB_ID)')
do
  gcloud beta dataflow jobs cancel $id
done

echo "Deleting BigQuery dataset and tables..."
bq rm -r -f tfeel

echo "Deleting PubSub topics and subscriptions..."
gcloud beta pubsub subscriptions delete $PUBSUB_TW_SUB
gcloud beta pubsub topics delete $PUBSUB_TW_TOPIC
gcloud beta pubsub topics delete $PUBSUB_RZ_TOPIC

echo "Delete Service Account, will prompt..."
gcloud beta iam service-accounts delete $SERVICE_ACCOUNT_NAME
