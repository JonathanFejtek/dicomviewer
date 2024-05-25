FROM golang:1.22.2

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN chmod +x entrypoint.sh

RUN go build ./cmd/dicomviewer

ENTRYPOINT ["./entrypoint.sh"]