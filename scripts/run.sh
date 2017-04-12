#!/bin/bash

# get global vars
. scripts/config.sh

export GOOGLE_APPLICATION_CREDENTIALS="./service-account-key.json"

../tfeel/tfeel --query="google, gcp, bigquery, cloud"
