FROM some-proget-server
LABEL authors="stoutes"

WORKDIR /app
COPY . .
RUN go test ./...
RUN go build ./cmd/server