docker run -d \
  -it \
  --env-file vals.txt \
  --restart=always \
  --name discordstonk \
  quay.io/nibalizer/discordstonksbot
