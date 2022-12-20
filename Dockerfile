FROM golang:alpine as builder
WORKDIR /build
COPY go.mod go.sum /docs/ya-root.crt static /docs/gDriveCredentials.json ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main main.go
FROM scratch
COPY --from=builder main /bin/main
ADD static /static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/ya-root.crt /etc/ssl/certs/
COPY --from=builder /build/gDriveCredentials.json /etc/ssl/certs
EXPOSE 3000 80
ENTRYPOINT ["/bin/main"]