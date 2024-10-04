FROM node:20.12-alpine as node-builder
    WORKDIR /x
    RUN npm install -g pnpm
    COPY . .
    RUN pnpm install
    RUN pnpm build

FROM golang:1.22-alpine as go-builder
    WORKDIR /x
    COPY . .
    COPY --from=node-builder /x/static/dist ./static/dist
    RUN go mod download
    RUN go mod verify
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /x/bin/atc ./cmd/atc/main.go

FROM gcr.io/distroless/static-debian11
    COPY --from=go-builder /x/bin/atc /
    CMD ["/atc"]
