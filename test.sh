#!/bin/bash
docker stop node-manager || true &&
docker rm node-manager || true &&
docker rmi node-manager || true &&
./build.sh &&
docker run -d -p 80:80 --name node-manager woz5999/node-manager &&
sleep 2 &&
docker logs node-manager
