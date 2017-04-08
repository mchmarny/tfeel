#!/bin/bash

# Assumes Twitter API envirnment variables are arleady defined
# export T_CONSUMER_KEY=""
# export T_CONSUMER_SECRET=""
# export T_ACCESS_TOKEN=""
# export T_ACCESS_SECRET=""

# Google
export APP_NAME=$(basename "$PWD")
export PROJECT_ID="mchmarny-dev"
export PUBSUB_TW_TOPIC="tweets"
export PUBSUB_RZ_TOPIC="results"
export PUBSUB_TW_SUB="${PUBSUB_TW_TOPIC}-events"
export SERVICE_ACCOUNT_NAME="${APP_NAME}-user"
export SA_EMAIL="${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

# ==============================================================================
# Service Account
# ==============================================================================

echo "Checking if Service Account alredy created..."
SA=$(gcloud beta iam service-accounts list --format='value(NAME)' --filter="NAME:${SERVICE_ACCOUNT_NAME}")
if [ -z "${SA}" ]; then
  echo "Service Account not set, creating..."
  gcloud beta iam service-accounts create ${SERVICE_ACCOUNT_NAME} \
    --display-name "${APP_NAME} app user"
fi

echo "Creating service account bindings..."
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${SA_EMAIL}" \
    --role='roles/pubsub.publisher' \
    --role='roles/pubsub.subscriber' \
    --role='roles/dataflow.developer' \
    --role='roles/dataflow.worker'

echo "Creating service account key..."
gcloud beta iam service-accounts keys create --iam-account $SA_EMAIL service-account-key.json

# ==============================================================================
# GCP Dependancies
# ==============================================================================

echo "Creating topics and subscriptions..."
gcloud beta pubsub topics create ${PUBSUB_TW_TOPIC}
gcloud beta pubsub topics create ${PUBSUB_RZ_TOPIC}
gcloud beta pubsub subscriptions create $PUBSUB_TW_SUB --topic=${PUBSUB_TW_TOPIC}

echo "Creating BigQuery tables..."
bq mk tfeel
bq mk --schema query:string,id:string,on:string,by:string,body:string -t tfeel.tweets
bq mk --schema id:string,score:float,parts:string -t tfeel.results

# projects/mchmarny-dev/topics/tweets
# mchmarny-dev:tfeel.tweets

echo "Creating dataflow job to drain tweet topic to BigQuery..."
gcloud beta dataflow jobs run ${APP_NAME}-${PUBSUB_TW_TOPIC} \
  --gcs-location gs://dataflow-templates/pubsub-to-bigquery/template_file \
  --parameters="topic=projects/${PROJECT_ID}/topics/${PUBSUB_TW_TOPIC}","table=${PROJECT_ID}:${APP_NAME}.${PUBSUB_TW_TOPIC}"

echo "Creating dataflow job to drain result topic to BigQuery..."
gcloud beta dataflow jobs run "${APP_NAME}-${PUBSUB_RZ_TOPIC}" \
  --gcs-location gs://dataflow-templates/pubsub-to-bigquery/template_file \
  --parameters="topic=projects/${PROJECT_ID}/topics/${PUBSUB_RZ_TOPIC}","table=${PROJECT_ID}:${APP_NAME}.${PUBSUB_RZ_TOPIC}"
