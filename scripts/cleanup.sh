#!/bin/bash

DIR="$(dirname "$0")"
. "${DIR}/config.sh"


read -p "Are you sure you want to cancel all Dataflow jobs created in this project today? [Y/n] " -n 1 -r
echo    # move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
  echo "Canceling all Dataflow jobs created today in this project..."
  TODAY=$(date +'%Y-%m-%d')
  for JOB_ID in $(gcloud beta dataflow jobs list --format='value(JOB_ID)')
  do
    JOB_DATE=$(echo ${JOB_ID}| cut -d'_' -f 1)
    if [ "$TODAY" == "$JOB_DATE" ]
    then
      gcloud beta dataflow jobs cancel $JOB_ID
    fi
  done
fi

echo "Deleting BigQuery dataset and tables..."
bq rm -r -f tfeel

echo "Deleting PubSub topics and subscriptions..."
gcloud beta pubsub subscriptions delete $PUBSUB_TW_SUB
gcloud beta pubsub topics delete $PUBSUB_TW_TOPIC
gcloud beta pubsub topics delete $PUBSUB_RZ_TOPIC


echo "Delete Service Account project bindings, will prompt..."
# TODO: why this fails due to insufficient caller permissions?
#gcloud beta iam service-accounts delete $SERVICE_ACCOUNT_NAME
gcloud projects remove-iam-policy-binding ${GCLOUD_PROJECT} \
    --member "serviceAccount:${SA_EMAIL}" --role roles/editor
