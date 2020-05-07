
# Build stage
FROM golang:1.13
ENV ROOT=/webserver
ADD . $ROOT
WORKDIR $ROOT
RUN make build
CMD ["./bin/webserver"]