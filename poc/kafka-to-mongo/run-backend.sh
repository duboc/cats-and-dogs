#!/bin/bash

docker run -d -p 9090:9090 -e KAFKA_URL=$WF_PROXY -e MONGO_URL="$MONGO_URL" duboc/catdog-kafka-to-mongo:1.0
