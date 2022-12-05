FROM golang:alpine as builder
WORKDIR /build
COPY go.mod go.sum static ./
# COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main main.go
FROM scratch
COPY --from=builder main /bin/main
ADD static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 3000 80
ENTRYPOINT ["/bin/main"]