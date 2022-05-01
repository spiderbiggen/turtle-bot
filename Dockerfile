# syntax=docker/dockerfile:1

FROM golang:1.18 as Builder

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./main.go ./main.go
COPY ./src ./src

RUN CGO_ENABLED=0 go build -o /weeb_bot

FROM gcr.io/distroless/static as Application

WORKDIR /opt

COPY --from=Builder /weeb_bot /weeb_bot

EXPOSE 8080
ENTRYPOINT ["/weeb_bot"]