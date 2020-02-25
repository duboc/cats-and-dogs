#!/bin/bash

docker run --restart=always -d -e NATS_ENDPOINT=$NATS_ENDPOINT -e MONGODB_ENDPOINT=$MONGODB_ENDPOINT duboc/cdbackend-mongo:1.0
