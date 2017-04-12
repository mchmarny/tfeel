#!/bin/bash

# get global vars
. scripts/config.sh

# create cluster
gcloud container --project ${GCLOUD_PROJECT} \
  clusters create "${APP_NAME}-cluster" \
  --zone "us-west1-b" \
  --machine-type "n1-standard-2"
  --image-type "COS" \
  --disk-size "100"
  --scopes "https://www.googleapis.com/auth/compute",\
    "https://www.googleapis.com/auth/devstorage.read_only",\
    "https://www.googleapis.com/auth/logging.write",\
    "https://www.googleapis.com/auth/pubsub",\
    "https://www.googleapis.com/auth/trace.append" \
  --num-nodes "3" \
  --network "default" \
  --enable-cloud-logging \
  --no-enable-cloud-monitoring
