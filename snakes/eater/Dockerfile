FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY . .

RUN go build -o eater .

FROM alpine

WORKDIR /opt
COPY --from=builder /build/eater .

EXPOSE 8080

CMD ["./eater"]


