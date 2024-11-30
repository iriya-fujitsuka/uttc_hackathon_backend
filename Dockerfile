# Build Stage
FROM golang:1.21 AS build

WORKDIR /app

COPY main.go .

# Initialize go.mod
RUN go mod init app

# スタティックリンクされた Linux 用バイナリを生成
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o app main.go

# Run Stage
FROM gcr.io/distroless/static

COPY --from=build /app/app /

EXPOSE 8080

CMD ["/app"]
