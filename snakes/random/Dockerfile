FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY . .

RUN go build -o random .

FROM alpine

WORKDIR /opt
COPY --from=builder /build/random .

EXPOSE 8080

CMD ["./random"]


