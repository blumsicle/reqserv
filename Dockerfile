FROM --platform=$BUILDPLATFORM golang:latest AS build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY .git go.mod go.sum Makefile ./
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} make deps

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build

FROM alpine:latest

COPY --from=build /app/bin/* /app
EXPOSE 8080
ENTRYPOINT ["/app", "run", "--config", "/config.yaml"]
