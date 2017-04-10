#!/bin/bash

dir="$(dirname "$0")"
. "${dir}/config.sh"

# ==============================================================================
# Service Account
# ==============================================================================

echo "Checking if Service Account alredy created..."
SA=$(gcloud beta iam service-accounts list --format='value(EMAIL)' --filter="EMAIL:${SA_EMAIL}")
if [ -z "${SA}" ]; then
  echo "Service Account not set, creating..."
  gcloud beta iam service-accounts create ${SERVICE_ACCOUNT_NAME} \
    --display-name="${APP_NAME} service account"
fi

# Set policies manualy
# gcloud iam service-accounts get-iam-policy ${SA_EMAIL} --format json > policy.json
# gcloud iam service-accounts set-iam-policy ${SA_EMAIL} policy.json

echo "Creating service account bindings..."
gcloud projects add-iam-policy-binding ${GCLOUD_PROJECT} \
    --member "serviceAccount:${SA_EMAIL}" --role roles/editor

echo "Creating service account key..."
gcloud beta iam service-accounts keys create --iam-account $SA_EMAIL \
  service-account-key.json
