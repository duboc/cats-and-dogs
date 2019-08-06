#!/bin/bash 

for ((i=1;i<=10000000;i++)); do   curl  --header "Connection: keep-alive" "http://ec2-54-89-29-170.compute-1.amazonaws.com:9090/api/dog" & done;

for ((i=1;i<=109990;i++)); do   curl  --header "Connection: keep-alive" "http://ec2-54-89-29-170.compute-1.amazonaws.com:9090/api/cat" & done;

for ((i=1;i<=10000;i++)); do   curl "http://ec2-54-89-29-170.compute-1.amazonaws.com" & done;
