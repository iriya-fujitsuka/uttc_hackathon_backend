FROM --platform=linux/amd64 golang:1.23.3 as build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/main ./main.go

# 環境変数の設定（必要に応じて変更）
ENV MYSQL_ROOT_PASSWORD=root_password \
    MYSQL_DATABASE=hackathon \
    MYSQL_USER=uttc \
    MYSQL_PASSWORD=1234 \
    MYSQL_HOST=unix(/cloudsql/term6-iriya-fujitsuka:us-central1:uttc)

FROM gcr.io/distroless/base
WORKDIR /root
COPY --from=build /app/cmd/main .
EXPOSE 8080
CMD ["./main"]