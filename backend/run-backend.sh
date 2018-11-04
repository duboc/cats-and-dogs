#!/bin/bash

docker run -d -p 9090:9090 -e WF_PROXY=$WF_PROXY duboc/cdbackend:1.1
