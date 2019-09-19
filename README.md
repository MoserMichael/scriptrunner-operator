# k8s operator that runs python scripts on pods

This is an exercise to learn about kubernetes, kubernetes operators and the operator sdk;

This operator defines a ScriptRunner object type, this object type is defined by a CRD (custom resource descriptor) that is part of the operator.

The Spec part of the ScriptRunner object defines the following:

1. the code of a python3 script, this script will be run by a pod running one each one of the selected nodes.
2. specify the set of nodes where the script will be run (by means of label selector)
3. If the script is to be run in either normal mode, privileged mode - a docker with elevated privileges, privileged with file system - an elevated docker with a mount to the hosts file system.
3. number of times that the script will be run
4. size of the command history (how many results to keep in the status section of the CRD per node runner)
5. optionally define stuff to be installed on the pod runner: pip modules to install and stuff to install via dnf install.

For each invocation of the script on a node you will get a entry in the status of the CRD object; the entry will list the result of the script invocation per node.

Here are some notes i took while learning about [k8s operators](https://github.com/MoserMichael/cstuff/releases/download/k8sop/kubernetes-operators.pdf)

The code
========

The code for this project has a commit for each stage of working with the operator-sdk, this way you can easily examine and distinguish generated code from hand written code.
Actually this is a confusing detail that has to be figured out when working in an environment that uses code generation: often it is important to discern the difference between generated code end and hand written code.
I think there is a general trick when working in an environment like this: keep a git repository that records each step of the process.

Build
=====

there is a ./build.sh script in the root direcory of the root directory.

Prerequisites are:
1. docker must be installed and runnable from the current user (not root).
2. golang must be installed
3. operator-sdk must be installed
4. minikube must be installed (did the testing with minikube) (preferable kvm2 must be installed along with minikube)

Note: always check version of kubernetes cluster is according to the operator-sdk requirements. (I did burn some time on this one)

The build script is a bash script that does the following steps:

1. set docker environment to local minikube repository: eval $(minikube docker-env)
2. build operator sdk and put it into docker image: scriptrunner:v0.0.1
3. tag that image to scriptrunner:latest (for whatever reason it only worked if there is a latest tag on the docker image)
4. build the docker image that runs the script on a k8s node in docker image scriptrunnerpod:v0.0.1 ; also tag it as latest
5. register the operator CRD, operators service account, role, role binding and deployment objects.
6. grant permissive binding, for whatever reason without this step the operator can't enumerate the nodes set. $(kubectl create clusterrolebinding permissive-binding   --clusterrole=cluster-admin   --user=admin   --user=kubelet   --group=system:serviceaccounts)

There is a ./test.sh script that runs some integration test. (Sorry, no unit testing in this project)

What i learned
==============

When running an elevated Pod its capabilities are determined by the environment that runs the pod: With minikube I would run the cluster in a kvm2 VM, in this configuration the pod sees the file system of the VM; when running the cluster in a container then we get the actual host file system. Fun.

Issues
======

1. Can't define the script as a multiline yamls field in the CRD; instead ; Somehow when I used multiline yaml messages in the instance definition we get a message that fails to reach the operator-sdk reconcile method, and we get the following error
<code>
E0918 14:27:57.855078       1 reflector.go:134] pkg/mod/k8s.io/client-go@v0.0.0-20190228174230-b40b2a5939e4/tools/cache/reflector.go:95: Failed to list \*unstructured.Unstructured: the server could not find the requested resource
E0918 14:27:58.861843       1 reflector.go:134] pkg/mod/k8s.io/client-go@v0.0.0-20190228174230-b40b2a5939e4/tools/cache/reflector.go:95: Failed to list \*v1alpha1.ScriptRunner: the server could not find the requested resource (get scriptrunners.scriptrunner.github.com)
E0918 14:27:58.863071       1 reflector.go:134] pkg/mod/k8s.io/client-go@v0.0.0-20190228174230-b40b2a5939e4/tools/cache/reflector.go:95: Failed to list \*unstructured.Unstructured: the server could not find the requested resource
E0918 14:28:03.809223       1 streamwatcher.go:109] Unable to decode an event from the watch stream: unable to decode watch event: v1alpha1.ScriptRunner.Spec: v1alpha1.ScriptRunnerSpec.PythonScript: []string: decode slice: expect [ or n, but found ", error found in #10 byte of ...|nScript":"print('hel|..., bigger context ...|oInstall":"","podType":"elevated","pythonScript":"print('hello operator')","runPeriod":1}}|...
</code>

2. The operator code needs to enumerate the node set, so that the pods can be directly scheduled to these nodes; It turns out that enumerating nodes is a privileged operation; I should rather have used pod deployment object, or something of that sort. Kubernetes doesn't like it when one is reinventing the wheel...

3. The current operator-sdk version comes with a controller and controller client that don't support the patch method for generated CRD object; had to use the client-go untyped client to perform a patch operation (by pod running script that modifiesthe CRD object status fields)

4. the current operator-sdk build command does not report compilation errors, instead it reports the command line that is used to build the operator code. No problem - the compile.sh script in this repository parses the error message and runs the compilation step. Still I think that this detail is a bit annoying.

