FROM golang:1.25.3-alpine3.22 AS build

COPY . /app
WORKDIR /app
RUN go build -o tcprcon

FROM alpine:3.22 AS run
COPY --from=build /app/tcprcon /app/tcprcon
WORKDIR /app
ENTRYPOINT [ "./tcprcon" ]
