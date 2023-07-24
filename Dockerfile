# syntax=docker/dockerfile:1

FROM golang:1.20 as Builder

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./main.go ./main.go
COPY ./internal ./internal

RUN CGO_ENABLED=0 go build -o /turtle-bot

FROM gcr.io/distroless/static as Application

WORKDIR /opt

COPY --from=Builder /turtle-bot /turtle-bot

EXPOSE 8080
ENTRYPOINT ["/turtle-bot", "-level=info"]