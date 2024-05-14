FROM golang:1.22.2 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o server cmd/server/main.go

#RUN apt-get update && apt-get install -y tesseract-ocr tesseract-ocr-eng && rm -rf /var/lib/apt/lists/*

FROM scratch
COPY --from=builder /app/server /server

COPY /config/config.yaml /config/config.yaml

COPY /migrations/deploy /migrations/deploy

CMD ["/server"]
