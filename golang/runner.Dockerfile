# ---------- ビルドフェーズ ----------
FROM golang:1.25.1-trixie AS builder

WORKDIR /app/web-app
COPY web-app/ ./

RUN go mod init web-app && go build -o server main.go

# ---------- 実行フェーズ ----------

FROM docker:28.5.0

WORKDIR /app

# install openldap-clients
RUN apk add --no-cache openldap-clients

# ビルド成果物をコピー
COPY --from=builder /app/web-app/server .

# ポート公開
EXPOSE 8082

# アプリケーション実行
CMD ["./server"]
