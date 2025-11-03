# ---------- ビルドフェーズ ----------
FROM golang:1.25.1-trixie AS builder

WORKDIR /app
COPY main.go .

RUN go mod init ldap-init && go build -o exec main.go

# ---------- 実行フェーズ ----------
FROM cleanstart/openldap:2.6-amd64

# 独立ディレクトリ作成とバイナリ配置
COPY --from=builder /app/exec /customized/exec

# ENTRYPOINT は Goバイナリのみ
ENTRYPOINT ["/customized/exec"]