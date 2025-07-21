# Dockerfile for ResponsibleAPI
FROM golang:1.24-alpine

# RUN go get github.com/githubnemo/CompileDaemon
RUN go install -mod=mod github.com/githubnemo/CompileDaemon
ENTRYPOINT ["CompileDaemon", "--build=go build -o ./bin/api ./cmd/api", "--command=./bin/api", "--polling"]

EXPOSE 8080