# to build this docker image:
#   docker build .
FROM ghcr.io/hybridgroup/opencv:4.8.1

ENV GOPATH /go

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

RUN go mod download 
RUN go mod verify

COPY . .
RUN go build -o main ./cmd/main.go

RUN chmod +x /app/entrypoint.sh

EXPOSE ${APP_PORT}

ENTRYPOINT ["/app/entrypoint.sh"]