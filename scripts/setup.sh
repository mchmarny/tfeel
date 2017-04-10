#!/bin/bash

DIR="$(dirname "$0")"
. "${DIR}/config.sh"

# ==============================================================================
# Service Account
# ==============================================================================

. "${DIR}/sa.sh"

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
