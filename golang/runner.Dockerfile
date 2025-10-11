FROM docker:28.5.0

WORKDIR /app

# ビルド成果物をコピー
COPY server .

# ポート公開
EXPOSE 8082

# アプリケーション実行
CMD ["./server"]
