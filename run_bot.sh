#!/bin/bash

# build command:
# docker build --no-cache -t quay.io/nibalizer/discordstonksbot:latest

docker pull quay.io/nibalizer/discordstonksbot
docker stop discordstonk
docker rm discordstonk

docker run -d \
  -it \
  --env-file vals.txt \
  --restart=always \
  --name discordstonk \
  quay.io/nibalizer/discordstonksbot:latest
