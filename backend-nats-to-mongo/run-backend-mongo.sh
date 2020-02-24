#!/bin/bash

docker run --restart=always -d -p 27017:27017 -e NATS_ENDPOINT= -e MONGODB_ENDPOINT= duboc/cdbackend-mongo:1.0
