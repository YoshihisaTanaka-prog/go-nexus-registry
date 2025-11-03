# ---------- ビルドフェーズ ----------
FROM golang:1.25.1-trixie AS builder

ARG TARGET_NAME

WORKDIR /app
COPY build.sh .
COPY ${TARGET_NAME}/ ./${TARGET_NAME}/

RUN ./build.sh ${TARGET_NAME}

# ---------- 実行フェーズ ----------

FROM debian:12

ARG TARGET_NAME

WORKDIR /app

# install openldap-clients
RUN (apt-get update && apt-get install -y ca-certificates curl gnupg lsb-release ldap-utils) && \ 
    (curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg) &&\
    (echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null) && \
    (apt-get update && apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin)

# ビルド成果物をコピー
COPY --from=builder /app/${TARGET_NAME}/server /app/server

# ポート公開
EXPOSE 8080

# アプリケーション実行
ENTRYPOINT ["/app/server"]
