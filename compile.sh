#!/bin/bash

IMG=$1

if [ "x$IMG" == "x" ]; then

    cat <<EOF
    ./compile <docker-image-name>

runs operator-sdk build command

if a compilation error occurs then current verion of operator-sdk tool show a message of the following form:
Error: failed to build operator binary: (failed to exec []string{"go", "build", ....
this tool parses this error message and runs the go build comman that it includes. very handy.
EOF
    exit 1
fi

unset GOROOT

eval $(minikube docker-env)

out=`operator-sdk build $IMG 2>&1`
STAT=$?

echo "$out"

if [ "x$STAT" != "x0" ]; then
    echo "operator-sdk failed: $STAT"

    if [[ $out =~ string\{([^\}]*) ]]; then
        CMP_LINE=${BASH_REMATCH[1]}
        #echo "cmp line: $CMP_LINE"
    fi

    CMD_LINE=$(echo "$CMP_LINE" | sed -e 's/", "/ /g' | sed 's/^.//' | sed 's/.$//')

    cat <<EOF

===
compile $CMD_LINE
====

EOF
    
    `$CMD_LINE`

fi    


exit $STAT


