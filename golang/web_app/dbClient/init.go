package dbClient

import (
	"context"
	"fmt"
	"os"
	"time"
	"web_app/ent"
	"web_app/customError"
	_ "github.com/lib/pq" // PostgreSQLドライバ
	_ "github.com/mattn/go-sqlite3" // SQLiteドライバ
)

func initDbUnit(client **ent.Client, dbType string, dsn string) {
	var err error
	// DBに接続
	*client, err = ent.Open(dbType, dsn)
	if err != nil {
		customError.Exit1(dbType, "への接続失敗: ", err)
	}

	// スキーマの自動マイグレーション（テーブルが無ければ作成）
	if err := (*client).Schema.Create(*ctx); err != nil {
		customError.Exit1(dbType, "スキーマ作成中のエラー: ", err)
	}
}

func InitDb(c *context.Context) {
	ctx = c
	time.Sleep(time.Second * 1)
	fmt.Fprintln(os.Stdout, "DBのスキーマを設定します。")

	initDbUnit(&psqlClient, "postgres", psqlDsn)
	initDbUnit(&ramClient, "sqlite3", ramDsn)

	fmt.Fprintln(os.Stdout, "DBのスキーマを設定しました。")
}