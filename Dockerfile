FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate

RUN go build -o ai-tracker

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/ai-tracker .

COPY --from=builder /app/website/static ./website/static

EXPOSE 8080

CMD ["./ai-tracker"]
