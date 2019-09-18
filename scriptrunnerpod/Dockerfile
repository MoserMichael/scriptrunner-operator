FROM fedora

RUN echo "get packages"
RUN dnf upgrade -y && dnf clean all -y && dnf install -y sudo golang python37 python3-pip

ADD ./pipinstall.sh /root/pipinstall.sh
ADD ./testcmdrun.sh /root/testcmdrun.sh

# works when host and image are the same. (problem docker build doesn't see network)
ADD ./scriptrunnerpod /root/commandrunner

# building it in the container environment.
#RUN echo "copy sources"
#ADD go.txt /root/go.mod
#ADD ./commandrunner.go /root/commandrunner.go
#RUN echo "build the package & get go packages"
#RUN cd /root && pwd >>/tmp/log.log && ls >>/tmp/log.log && go build ./commandrunner.go

#CMD [ "./commandrunner" ]
