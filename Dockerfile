FROM golang:1.18 as builder

WORKDIR /go/ports

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ports

RUN go test ./...

FROM alpine:3.15 as runner

COPY --from=builder /go/ports/ports .

#I'll mount this as a volume if I have time
COPY --from=builder /go/ports/ports.json .

ENTRYPOINT ./ports