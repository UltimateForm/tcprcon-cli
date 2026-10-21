FROM golang:1.26.10-alpine3.24 AS build

COPY . /app
WORKDIR /app
RUN go build -o tcprcon

FROM alpine:3.24 AS run
COPY --from=build /app/tcprcon /app/tcprcon
WORKDIR /app
ENTRYPOINT [ "./tcprcon" ]
