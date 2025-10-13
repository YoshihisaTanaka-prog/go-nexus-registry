package dbClient

import (
	"context"
	"fmt"
	"os"
	"time"
	"web_app/ent"
	"web_app/customError"
	_ "github.com/lib/pq" // PostgreSQLドライバ
	// _ "github.com/mattn/go-sqlite3" // SQLiteドライバ
)

func InitDb() {
	time.Sleep(time.Second * 1)
	fmt.Fprintln(os.Stdout, "DBのスキーマを設定します。")
	var err error
	client, err = ent.Open("postgres", dsn)
	if err != nil {
		customError.Exit1("PostgreSQLへの接続失敗: ", err)
	}

	// スキーマの自動マイグレーション（テーブルが無ければ作成）
	if err := client.Schema.Create(context.Background()); err != nil {
		customError.Exit1("スキーマ作成中のエラー: ", err)
	}
	fmt.Fprintln(os.Stdout, "DBのスキーマを設定しました。")
}