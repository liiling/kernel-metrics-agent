FROM golang:1.14.4
ENV GO111MODULE=on

# workdir in the image file system
WORKDIR /usr/src
COPY go.mod .
copy go.sum .
RUN go mod download
COPY . .