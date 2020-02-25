#!/bin/bash

docker run --restart=always -d -p 80:80  -e WAVEFRONT_INSTANCE=$WAVEFRONT_INSTANCE -e WAVEFRONT_TOKEN=$WAVEFRONT_TOKEN -e NATS_ENDPOINT=$NATS_ENDPOINT duboc/cdbackend-nats:1.0
