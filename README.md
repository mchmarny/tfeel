# TFEEL

Short for Twitter Feeling... simple sentiment analyses over tweeter data for
specific terms using Google Cloud services.

![tfeel data flow](https://raw.github.com/mchmarny/tfeel/master/images/tfeel-flow-chart.png)

> All GCP services used in this example can be run under the GCP Free Tier
plan. More more information see https://cloud.google.com/free/

## Configuration

Edit the `scripts/config.sh` file with your Twitter API info. Alternatively
define the following environment variables

```
# export T_CONSUMER_KEY=""
# export T_CONSUMER_SECRET=""
# export T_ACCESS_TOKEN=""
# export T_ACCESS_SECRET=""
```

## Dependencies

`tfeel` depends on the following GCP services:

* [Cloud Pub/Sub](https://cloud.google.com/pubsub/)
* [BigQuery](https://cloud.google.com/bigquery/)
* [Dataflow](https://cloud.google.com/dataflow/)

You can set all these resource dependencies using the `scripts/setup.sh` script. This script will also configure a service account for the `tfeel`  application to run under.

## Run

Once all the necessary GC resources have been created (dependencies) you can execute the `tfeel` application using the `scripts/run.sh` script.

## Cleanup

The cleanup of all the resources created in this application can be accomplished by executing the `scripts/cleanup.sh` script.

### TODO

* Tests, yes please
* Minimize account service roles (currently, project editor)
* Add app containerization => [GKE](https://cloud.google.com/container-engine/)
* Add Spanner DB to keep user state (trend user-level sentiment)
