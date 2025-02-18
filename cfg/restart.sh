#!/usr/bin/env bash

sudo pm2 stop hls-streamer
sudo GOMAXPROCS=3 pm2 start streamer-linux-amd64 --name=hls-streamer -- -config=./prod.ini
sudo pm2 save