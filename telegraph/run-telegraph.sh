#!/bin/bash

docker run -v $PWD/telegraf.conf:/etc/telegraf/telegraf.conf:ro telegraf