
FROM golang:1.22.4 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y gcc
ENV CGO_ENABLED=1

RUN go build -o forum ./main.go

RUN ls -l forum

FROM golang:1.22.4

COPY --from=builder /app/forum /app/forum

COPY --from=builder /app/handlers /app/handlers
COPY --from=builder /app/web /app/web
COPY --from=builder /app/uploads /app/uploads

WORKDIR /app

CMD ["./forum"]