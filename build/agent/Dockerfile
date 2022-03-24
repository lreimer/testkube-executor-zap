# syntax=docker/dockerfile:1
FROM golang:1.17 as builder

WORKDIR /build
COPY . .
ENV CGO_ENABLED=0 
ENV GOOS=linux

RUN cd cmd/agent;go build -o /runner -mod mod -a .

FROM owasp/zap2docker-stable:2.11.1
COPY --from=builder /runner /zap/runner
ENTRYPOINT ["/zap/runner"]
