FROM golang:latest as build-env
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main
FROM scratch
COPY --from=build-env /app/main /app/main
ENTRYPOINT ["/app/main"]
