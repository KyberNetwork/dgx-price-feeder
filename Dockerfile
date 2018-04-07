# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/KyberNetwork/dgx-price-feeder

WORKDIR /go/src/github.com/KyberNetwork/dgx-price-feeder
RUN go install -v github.com/KyberNetwork/dgx-price-feeder/cmd

EXPOSE 8000

