#!/bin/bash

docker pull quay.io/nibalizer/discordstonksbot
docker stop discordstonk
docker rm discordstonk

docker run -d \
  -it \
  --env-file vals.txt \
  --restart=always \
  --name discordstonk \
  quay.io/nibalizer/discordstonksbot:latest
