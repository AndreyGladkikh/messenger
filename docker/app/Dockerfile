# syntax=docker/dockerfile:1

FROM golang:1.26 AS build-dev

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum* ./
RUN go mod download

COPY ./ ./
COPY ./configs /etc/app/configs

RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /bin/app ./cmd/web/main.go

EXPOSE 40000 40040 8080 50051

USER 1000:1000

ENTRYPOINT ["/go/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/bin/app"]

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
