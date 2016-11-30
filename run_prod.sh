#!/bin/bash

docker pull $1/nufito-prod:latest
if docker stop nufito-app; then docker rm nufito-app; fi
echo "running image"
docker run -d -p 8080:8080 --net=host --name nufito-app $1/nufito-prod
echo "running docker rmi"
if docker rmi $(docker images --filter "dangling=true" -q --no-trunc); then :; fi
