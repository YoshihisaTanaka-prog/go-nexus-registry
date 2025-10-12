# ---------- ビルドフェーズ ----------
FROM golang:1.25.1-trixie AS builder

ARG TARGET_NAME

WORKDIR /app
COPY build.sh /app
COPY ${TARGET_NAME}/ ./${TARGET_NAME}/

RUN ./build.sh ${TARGET_NAME}

# ---------- 実行フェーズ ----------

FROM docker:28.5.0

ARG TARGET_NAME

WORKDIR /app

# install openldap-clients
RUN apk add --no-cache openldap-clients

# ビルド成果物をコピー
COPY --from=builder /app/${TARGET_NAME}/server .

# ポート公開
EXPOSE 8082

# アプリケーション実行
CMD ["./server"]
