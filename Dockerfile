FROM golang:latest

# copy local dir
COPY . /go/steward
WORKDIR /go/steward

RUN go install .
CMD ["sleep", "infinity"]