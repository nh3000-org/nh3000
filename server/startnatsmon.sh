#!/bin/sh
tailf /opt/nats/log | /opt/log -logpattern DBG -serverip nats://192.168.0.103:4222
