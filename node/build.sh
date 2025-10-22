#!/usr/bin/sh

if [ ! -d node_modules ]; then
  npm ci .
fi

npm run build; chmod -R 777 dist