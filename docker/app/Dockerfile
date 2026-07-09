# syntax=docker/dockerfile:1

FROM golang:1.26 AS build-prod

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/web/main.go

FROM build-prod as build-test

COPY ./configs /etc/app/configs

EXPOSE 8080 50051

USER 1000:1000

ENTRYPOINT ["/bin/app"]

FROM gcr.io/distroless/base-debian11 AS build-release-prod

WORKDIR /

COPY --from=build-prod /app/configs /etc/app/configs
COPY --from=build-prod /bin/app /bin/app

EXPOSE 8080 50051

USER 1000:1000

ENTRYPOINT ["/bin/app"]
