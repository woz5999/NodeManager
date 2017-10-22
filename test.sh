#!/bin/bash
docker stop node-manager || true &&
docker rm node-manager || true &&
# docker rmi node-manager || true &&
# ./build.sh &&
export NODEMAN_AWS_REGION=us-east-1
export NODEMAN_AWS_SQS_QUEUE_URL='https://sqs.us-east-1.amazonaws.com/282028653949/scaling-test'
export NODEMAN_DEBUG=true
docker rmi node-manager || true &&
./build.sh &&
docker run -d -p 80:80 \
-e NODEMAN_AWS_REGION \
-e NODEMAN_AWS_SQS_QUEUE_URL \
-e NODEMAN_DEBUG \
-e AWS_ACCESS_KEY_ID \
-e AWS_SECRET_ACCESS_KEY \
--name node-manager woz5999/node-manager &&
sleep 2 &&
docker logs node-manager
