FROM nginx:1.29-alpine

# envsubst を利用するために gettext をインストール
RUN apk add --no-cache gettext

# 設定テンプレートをコピー
COPY templates/ /etc/nginx/

# エントリポイントスクリプトを追加
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh