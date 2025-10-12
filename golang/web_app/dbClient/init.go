package dbClient

import (
	"context"
	"fmt"
	"os"
	"web_app/ent"
	_ "github.com/lib/pq" // PostgreSQLドライバ
	// _ "github.com/mattn/go-sqlite3" // SQLiteドライバ
)

func InitDb() {
	var err error
	client, err = ent.Open("postgres", dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "PostgreSQLへの接続失敗: ", err)
	}

	// スキーマの自動マイグレーション（テーブルが無ければ作成）
	if err := client.Schema.Create(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "スキーマ作成中のエラー: ", err)
	}
}