#!/bin/sh

#ip-tables for nats with app
iptables -A INPUT -p tcp --destination-port 80 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT -p tcp --destination-port 4222 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT -p tcp --destination-port 8222 -m state --state NEW,ESTABLISHED -j ACCEPT

iptables -A OUTPUT  -p tcp --source-port 80 -m state --state ESTABLISHED -j ACCEPT
iptables -A OUTPUT  -p tcp --source-port 4222 -m state --state ESTABLISHED -j ACCEPT
iptables -A OUTPUT  -p tcp --source-port 8222 -m state --state ESTABLISHED -j ACCEPT
