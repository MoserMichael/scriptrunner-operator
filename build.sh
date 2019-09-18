#!/bin/bash

set -ex

minikube start --vm-driver=kvm2

eval $(minikube docker-env)

docker images

operator-sdk build scriptrunner:v0.0.1

# only managed to refer to images by latest tag from pod; strange.
docker tag scriptrunner:v0.0.1 scriptrunner:latest

pushd scriptrunnerpod

go build

docker build -t scriptrunnerpod:v0.0.1 .

# only managed to refer to images by latest tag from pod; strange.
docker tag scriptrunnerpod:v0.0.1 scriptrunnerpod:latest

popd


## register crd
kubectl create -f deploy/crds/scriptrunner_v1alpha1_scriptrunner_crd.yaml 
# register deployment/role/etc.
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml


#
# doing this on minikube:  magic command to create 'permissive binding': otherwise not allowed to list nodes in the controller.
#
kubectl create clusterrolebinding permissive-binding   --clusterrole=cluster-admin   --user=admin   --user=kubelet   --group=system:serviceaccounts


echo "** completed **"



