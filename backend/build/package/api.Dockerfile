# syntax=docker/dockerfile:1

FROM golang:1.27.1 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

FROM build-stage AS run-test-stage

RUN go test -v ./...

# FROM build-stage as build-test

# COPY ./configs /etc/app/configs

# EXPOSE 8080 50051

# USER 1000:1000

# ENTRYPOINT ["/bin/app"]

FROM gcr.io/distroless/base-debian13 AS build-release-stage

WORKDIR /

# COPY --from=build-stage /app/configs /etc/app/configs
COPY --from=build-stage /bin/api /bin/api

EXPOSE 8080 50051

USER nonroot:nonroot

ENTRYPOINT ["/bin/api"]
