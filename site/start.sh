#!/bin/bash

cat /usr/share/nginx/html/js/boot.js | envsubst '${BACKEND_URL}' > /usr/share/nginx/html/js/boot.js

nginx -g 'daemon off;'

