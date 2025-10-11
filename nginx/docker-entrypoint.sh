#!/bin/sh
set -e

# Nginx設定ファイルを環境変数で置換して生成
envsubst '${NEXUS_EXPOSED_HOST_NAME}' \
  < /etc/nginx/templates/default.conf.template \
  > /etc/nginx/conf.d/default.conf

echo "✅ Generated /etc/nginx/conf.d/default.conf with env variables:"
cat /etc/nginx/conf.d/default.conf
echo "----------------------------------------"

if [ -e /etc/nginx/templates/nexus_blocked_api_users.map.template ]; then
  cp /etc/nginx/templates/nexus_blocked_api_users.map.template /etc/nginx/conf.d/nexus_blocked_api_users.map
else
  echo "Please setup ./nginx/templates/nexus_blocked_api_users.map.template"
  echo "up while referring to './nginx/templates/nexus_blocked_api_users.map.template.template'"
  exit 1
fi

# Nginx起動
exec "$@"
