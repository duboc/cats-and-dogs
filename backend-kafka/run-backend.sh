#!/bin/bash

docker run -d -p 9090:9090 -e WF_PROXY=$WF_PROXY -e KAFKA_PROXY=$KAFKA_PROXY duboc/cdbackend-kafka:1.0
