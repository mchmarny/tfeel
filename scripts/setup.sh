#!/bin/bash

DIR="$(dirname "$0")"
. "${DIR}/config.sh"

# ==============================================================================
# Service Account
# ==============================================================================

# Set policies manualy
# gcloud iam service-accounts get-iam-policy ${SA_EMAIL} --format json > policy.json
# gcloud iam service-accounts set-iam-policy ${SA_EMAIL} policy.json

echo "Checking if Service Account alredy created..."
SA=$(gcloud beta iam service-accounts list --format='value(EMAIL)' --filter="EMAIL:${SA_EMAIL}")
if [ -z "${SA}" ]; then
  echo "Service Account not set, creating..."
  gcloud beta iam service-accounts create ${SERVICE_ACCOUNT_NAME} \
    --display-name="${APP_NAME} service account"

  echo "Creating service account key..."
  gcloud beta iam service-accounts keys create --iam-account $SA_EMAIL \
    service-account-key.json
fi

echo "Creating service account bindings..."
gcloud projects add-iam-policy-binding ${GCLOUD_PROJECT} \
    --member "serviceAccount:${SA_EMAIL}" --role roles/editor

# ==============================================================================
# GCP Dependancies
# ==============================================================================

echo "Creating topics and subscriptions..."
gcloud beta pubsub topics create ${PUBSUB_TW_TOPIC}
gcloud beta pubsub topics create ${PUBSUB_RZ_TOPIC}
gcloud beta pubsub subscriptions create $PUBSUB_TW_SUB \
  --topic=${PUBSUB_TW_TOPIC} \
  --ack-deadline=60

echo "Creating BigQuery tables..."
bq mk tfeel
bq mk --schema query:string,id:string,on:string,by:string,body:string -t tfeel.tweets
bq mk --schema id:string,score:float,parts:string -t tfeel.results

echo "Creating dataflow job to drain tweet topic to BigQuery..."
gcloud beta dataflow jobs run ${APP_NAME}-${PUBSUB_TW_TOPIC} \
  --gcs-location gs://dataflow-templates/pubsub-to-bigquery/template_file \
  --parameters="topic=projects/${GCLOUD_PROJECT}/topics/${PUBSUB_TW_TOPIC}","table=${GCLOUD_PROJECT}:${APP_NAME}.${PUBSUB_TW_TOPIC}" \
  --service-account-email="${SA_EMAIL}"

echo "Creating dataflow job to drain result topic to BigQuery..."
gcloud beta dataflow jobs run "${APP_NAME}-${PUBSUB_RZ_TOPIC}" \
  --gcs-location gs://dataflow-templates/pubsub-to-bigquery/template_file \
  --parameters="topic=projects/${GCLOUD_PROJECT}/topics/${PUBSUB_RZ_TOPIC}","table=${GCLOUD_PROJECT}:${APP_NAME}.${PUBSUB_RZ_TOPIC}" \
  --service-account-email="${SA_EMAIL}"
