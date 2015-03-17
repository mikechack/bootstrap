# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/MediaFusion/mf-connector

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get "github.com/streadway/amqp"
RUN go install github.com/MediaFusion/mf-connector

# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/mf-connector"]
CMD ["-http-port=8080"]


# Document that the service listens on port 8080.
EXPOSE 8080


some moe stuff
