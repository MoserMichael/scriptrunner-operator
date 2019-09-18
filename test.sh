#!/bin/bash

set -ex

function one_test()
{
echo "**** testing ${TEST_NAME} ****"

rm -f log-operator-${TEST_NAME}.log || true
rm -f log-pod-${TEST_NAME}.log || true

kubectl get deployment

CR_FILE="deploy/crds/scriptrunner_v1alpha1_scriptrunner_cr-${TEST_NAME}.yaml"

stat $CR_FILE

kubectl apply -f $CR_FILE

sleep 2

kubectl get scriptrunner 

kubectl get scriptrunner $TEST_NAME -o json

set +e
SCRIPTRUNNER_POD=$(kubectl get pods -o name | grep $TEST_NAME)
set -e

if [ "x${SCRIPTRUNNER_POD}" != "x" ]; then
    kubectl logs  ${SCRIPTRUNNER_POD} >log-pod-${TEST_NAME}.log
else
    echo "error: scriptrunner pod not running for $TEST_NAME !"
    exit 1
fi

kubectl delete scriptrunner $TEST_NAME

kubectl get pods

OPERATOR_POD=$(kubectl get pods -o name | grep scriptrunner-operator)

kubectl logs $OPERATOR_POD >log-operator-${TEST_NAME}.log
}

TEST_NAME="example-scriptrunner2"
one_test

TEST_NAME="example-scriptrunner-ps"
one_test

TEST_NAME="example-scriptrunner-ps2"
one_test

echo " ** all tests passed ** "
