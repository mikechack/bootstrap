# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get "github.com/streadway/amqp"
#RUN go get -u sqbu.com/MediaFusion/bootstrap

#RUN go get 192.168.59.3/MediaFusion/bootstrap
RUN go get github.com/mikechack/bootstrap
#RUN go install sqbu.com/MediaFusion/bootstrap


# Run the outyet command by default when the container starts.
#ENTRYPOINT ["/go/bin/bootstrap"]
#CMD ["-http-port=8080"]


# Document that the service listens on port 8080.
EXPOSE 8080
