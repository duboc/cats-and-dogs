#!/bin/bash

docker run --rm --name nginx -v $(pwd):/usr/share/nginx/html:ro -p 8080:80 nginx:1.15.0
