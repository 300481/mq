# mq

Abstraction for using message queues.

## Configuration

This module needs a configuration by environment variables listed below for each cloud provider.

### Google Cloud Platform PubSub

Environment Variable     |Description                                 |Type / valid values
-------------------------|--------------------------------------------|-----------------------
`GCP_CREDENTIALS_FILE`   |The path of the credentials file (json-file)|*String*
`GCP_TOPIC_NAME`         |The PubSub topic name                       |*String*
`GCP_CREATE_TOPIC`       |Allow creation of topic if not exists       |*String* `TRUE`/`FALSE`
`GCP_SUBSCRIPTION_NAME`  |The PubSub subscription name/id             |*String*
`GCP_CREATE_SUBSCRIPTION`|Allow creation of subscription if not exists|*String* `TRUE`/`FALSE`
`GCP_PROJECT_ID`         |The Google Project ID                       |*String*
