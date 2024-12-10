FROM golang:1.21.2-alpine3.18 AS build-stage
WORKDIR /notificationservice
COPY ./ /notificationservice
RUN mkdir -p /notificationservice/build
RUN go mod download
RUN go build -v -o /notificationservice/build/api ./cmd/server


FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /notificationservice/build/api /
COPY --from=build-stage /notificationservice/env/.env /
EXPOSE 3001
CMD ["/api"]