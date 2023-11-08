package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var Db *sql.DB

func Init() {
	// MySQLデータベースへの接続を初期化

	// ①-1: 環境変数からMySQLのユーザー名、パスワード、データベース名,ホスト名を取得
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	// ①-2: データベースへのDSN（Data Source Name）を構築
	connStr := fmt.Sprintf("%s:%s@%s/%s", mysqlUser, mysqlUserPwd, mysqlHost, mysqlDatabase)

	// ①-3: SQLデータベースに接続
	_db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Database connection error: %v\n", err)
	}
	if err = _db.Ping(); err != nil {
		log.Fatalf("Database ping error: %v\n", err)
	}
	Db = _db
}
