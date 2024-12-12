FROM golang:1.23.2-alpine AS build-stage
WORKDIR /notificationservice
COPY ./ /notificationservice
RUN mkdir -p /notificationservice/build
RUN go mod download
RUN go build -v -o /notificationservice/build/api ./cmd/server

# Use a minimal base image
FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /notificationservice/build/api /
COPY --from=build-stage /notificationservice/env/.env /etc/env/.env
EXPOSE 8012
CMD ["/api"]
