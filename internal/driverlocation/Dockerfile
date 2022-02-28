FROM golang:1.17.3-alpine3.14 as dev
RUN go get github.com/githubnemo/CompileDaemon
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait
WORKDIR /app/
COPY go.mod .
COPY ./pkg ./pkg
COPY ./internal/driverlocation ./internal/driverlocation
RUN go mod tidy
RUN go mod download
RUN go mod verify
WORKDIR /app/internal/driverlocation
CMD /wait && CompileDaemon -directory=. -directory=/app/pkg -build="go build -o driverlocation" -command="./driverlocation"
