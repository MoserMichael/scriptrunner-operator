#!/bin/bash

set -ex

go build

docker build -t scriptrunnerpod:v0.0.1 .

# only managed to refer to images by latest tag from pod; strange.
docker tag scriptrunnerpod:v0.0.1 scriptrunnerpod:latest


