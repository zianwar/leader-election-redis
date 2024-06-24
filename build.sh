#!/usr/bin/env bash
set -e

docker build -t server .
docker run -e PORT=8001 -p 8001:8001 server
