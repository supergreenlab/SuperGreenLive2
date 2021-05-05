#!/bin/bash

docker build -t liveserver-dev . -f Dockerfile.dev
docker run  --name=liveserver --network=supergreencloud_back-tier -p 8081:8081 --rm -it -v $(pwd)/config:/etc/liveserver -v $(pwd):/app liveserver-dev
