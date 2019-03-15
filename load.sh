#!/bin/bash 

loadtest -n 1000 -c 100 --rps 200 http://ec2-54-89-29-170.compute-1.amazonaws.com/ --quiet &
loadtest -n 10500 -c 100 --rps 200 http://ec2-54-89-29-170.compute-1.amazonaws.com:9090/api/cat --quiet &
loadtest -n 20100 -c 100 --rps 200  http://ec2-54-89-29-170.compute-1.amazonaws.com:9090/api/dog --quiet &
loadtest -n 100 http://ec2-54-89-29-170.compute-1.amazonaws.com/test --quiet &
