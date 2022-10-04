#!/bin/bash

docker build -t liveserver-dev . -f Dockerfile.dev
docker run --name=liveserver -p 8081:8081 --rm -it -e GOPRIVATE=github.com/SuperGreenLab/AppBackend -e SSH_AUTH_SOCK=/ssh-agent -v ${SSH_AUTH_SOCK}:/ssh-agent -v ~/.ssh/${git_github_identity:-id_rsa}:/root/.ssh/id_rsa -v $(pwd)/config:/etc/liveserver -v $(pwd):/app -v $(pwd)/storage:/tmp/storage liveserver-dev
