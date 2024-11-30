# Build Stage
FROM golang:1.21 AS build

WORKDIR /app

# go.modとgo.sumをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコード全体をコピー
COPY . .

# アプリケーションをビルド
RUN go build -o app main.go

# Run Stage
FROM gcr.io/distroless/static

COPY --from=build /app/app /

EXPOSE 8000

CMD ["/app"]