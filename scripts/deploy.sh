#!/bin/bash

DIR="$(dirname "$0")"
. "${DIR}/config.sh"


gcloud beta compute instances create "${APP_NAME}-vm" \
    --project=$GCLOUD_PROJECT \
    --zone=us-west1-c \
    --machine-type=n1-standard-2 \
    --scopes=cloud-platform \
    --boot-disk-size=10GB \
    --boot-disk-type=pd-ssd \
    --boot-disk-device-name=tfeel-vm


#gcloud compute ssh --project $GCLOUD_PROJECT --zone us-west1-c "${APP_NAME}-vm"