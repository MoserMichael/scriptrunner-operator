#!/bin/bash

set -ex

export GOROOT=$(dirname $(which go))

operator-sdk generate k8s

operator-sdk generate openapi

echo "** done **"
