FROM golang:1.22-alpine AS build

WORKDIR /src
COPY . ./
RUN go build -o /bin/service ./cmd/server

FROM alpine:3.19
COPY --from=build /bin/service /bin/service
EXPOSE 8080
CMD ["/bin/service"]
