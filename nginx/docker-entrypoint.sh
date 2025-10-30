#!/bin/sh
set -e

# Nginx設定ファイルを環境変数で置換して生成
envsubst '${NEXUS_EXPOSED_HOST_NAME}' \
  < /etc/nginx/templates/default.conf.template \
  > /etc/nginx/conf.d/default.conf

echo "✅ Generated /etc/nginx/conf.d/default.conf with env variables:"
echo ""
cat /etc/nginx/conf.d/default.conf
echo ""
echo "----------------------------------------"
echo ""

if [ -e /etc/nginx/templates/nexus_blocked_api_users.map.template ]; then
  cp /etc/nginx/templates/nexus_blocked_api_users.map.template /etc/nginx/conf.d/nexus_blocked_api_users.map
else
  echo "Please setup ./nginx/templates/nexus_blocked_api_users.map.template"
  echo "while referring to './nginx/templates/nexus_blocked_api_users.map.template.template'"
  exit 1
fi

sleep 10

set -e

entrypoint_log() {
  if [ -z "${NGINX_ENTRYPOINT_QUIET_LOGS:-}" ]; then
      echo "$@"
  fi
}

if [ "$1" = "nginx" ] || [ "$1" = "nginx-debug" ]; then
  if /usr/bin/find "/docker-entrypoint.d/" -mindepth 1 -maxdepth 1 -type f -print -quit 2>/dev/null | read v; then
    entrypoint_log "$0: /docker-entrypoint.d/ is not empty, will attempt to perform configuration"

    entrypoint_log "$0: Looking for shell scripts in /docker-entrypoint.d/"
    find "/docker-entrypoint.d/" -follow -type f -print | sort -V | while read -r f; do
      case "$f" in
        *.envsh)
          if [ -x "$f" ]; then
            entrypoint_log "$0: Sourcing $f";
            . "$f"
          else
            # warn on shell scripts without exec bit
            entrypoint_log "$0: Ignoring $f, not executable";
          fi
          ;;
        *.sh)
          if [ -x "$f" ]; then
            entrypoint_log "$0: Launching $f";
            "$f"
          else
            # warn on shell scripts without exec bit
            entrypoint_log "$0: Ignoring $f, not executable";
          fi
          ;;
        *) entrypoint_log "$0: Ignoring $f";;
      esac
    done

    entrypoint_log "$0: Configuration complete; ready for start up"
  else
    entrypoint_log "$0: No files found in /docker-entrypoint.d/, skipping configuration"
  fi
fi

# Nginx起動
exec "$@"
