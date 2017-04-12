#!/bin/bash

# Assumes Twitter API envirnment variables are arleady defined
# export T_CONSUMER_KEY=""
# export T_CONSUMER_SECRET=""
# export T_ACCESS_TOKEN=""
# export T_ACCESS_SECRET=""

# Google
export APP_NAME=$(basename "$PWD")
export PUBSUB_TW_TOPIC="tweets"
export PUBSUB_RZ_TOPIC="results"
export PUBSUB_TW_SUB="${PUBSUB_TW_TOPIC}-events"
export SERVICE_ACCOUNT_NAME="${APP_NAME}-sa"
export SA_EMAIL="${SERVICE_ACCOUNT_NAME}@${GCLOUD_PROJECT}.iam.gserviceaccount.com"
