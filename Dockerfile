# Build Stage
FROM golang:1.23.3 AS build

WORKDIR /app

# go.modとgo.sumをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコード全体をコピー
COPY . .

# Linux用のビルド設定
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# アプリケーションをビルド
RUN go build -o app main.go

# Run Stage
FROM gcr.io/distroless/static

COPY --from=build /app/app /

EXPOSE 8080

CMD ["/app/app"]