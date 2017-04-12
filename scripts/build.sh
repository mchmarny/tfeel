#!/bin/bash

# get global vars
. scripts/config.sh

docker build ./ \
  --build-arg T_CONSUMER_KEY="${T_CONSUMER_KEY}" \
  --build-arg T_CONSUMER_SECRET="${T_CONSUMER_SECRET}" \
  --build-arg T_ACCESS_TOKEN="${T_ACCESS_TOKEN}" \
  --build-arg T_ACCESS_SECRET="${T_ACCESS_SECRET}" \
  --build-arg GCLOUD_PROJECT="${GCLOUD_PROJECT}" \
  --build-arg APP_QUERY="google, gcp, bigquery, spanner"

# run
# docker run -i -t (+last IMAGE ID from 'docker imaages')
