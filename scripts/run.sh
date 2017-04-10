#!/bin/bash

# get global vars
. scripts/config.sh

export GOOGLE_APPLICATION_CREDENTIALS="./service-account-key.json"

# run
../tfeel/tfeel -q="google, gcp, bigquery, cloud" -p="${GCLOUD_PROJECT}"
