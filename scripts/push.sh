#!/bin/bash

# get global vars
. scripts/config.sh

# het image id (docker images)
LAST_IMAGE=$(docker images | grep -E '^golang.*latest' | awk '{print $3}')

# tag it
docker tag $LAST_IMAGE "gcr.io/${GCLOUD_PROJECT}/tfeel"

# push it
gcloud docker -- push "gcr.io/${GCLOUD_PROJECT}/tfeel:latest"
