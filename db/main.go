// main.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oklog/ulid"
)

type Item struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Category      string    `json:"category"`
	Chapter       string    `json:"chapter"`
	File          string    `json:"file"`
	CreatedBy     string    `json:"createdBy"`
	CreatedByName string    `json:"createdByName"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

var db *sql.DB

func init() {
	// MySQLデータベースへの接続を初期化

	// ①-1: 環境変数からMySQLのユーザー名、パスワード、データベース名を取得
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	// ①-2: データベースへのDSN（Data Source Name）を構築
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", mysqlUser, mysqlUserPwd, mysqlDatabase)

	// ①-3: SQLデータベースに接続
	_db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Database connection error: %v\n", err)
	}
	if err = _db.Ping(); err != nil {
		log.Fatalf("Database ping error: %v\n", err)
	}
	db = _db
}

func main() {
	// CORSミドルウェアを適用
	http.Handle("/api/addItem", CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 通常のリクエストの処理
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			// preflightリクエストの場合、200 OKを返して終了
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodPost:
			// POSTリクエストの処理
			handleAddItem(w, r)
		default:
			// サポートされていないメソッドの場合、405 Method Not Allowedを返す
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	log.Println("Server listening on :8000")
	http.ListenAndServe(":8000", nil)
}

// handleAddItem はPOSTリクエストを処理する関数
// 指定されたエラーメッセージをログに記録してHTTPエラーレスポンスを返すユーティリティ関数
func logAndSendError(w http.ResponseWriter, message string, status int, err error) {
	log.Printf("Error: %v\n", err)
	http.Error(w, message, status)
}

// handleAddItem はPOSTリクエストを処理する関数
func handleAddItem(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドがPOSTでない場合はエラーを返す
	if r.Method != http.MethodPost {
		logAndSendError(w, "Only POST requests are allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	// リクエストボディからデータをデコード
	var data Item
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logAndSendError(w, "Failed to decode request body", http.StatusBadRequest, err)
		return
	}

	// バリデーション: 必要なフィールドの欠落をチェック
	if data.Title == "" || data.Category == "" || data.Chapter == "" || data.CreatedBy == "" {
		logAndSendError(w, "Required fields are missing", http.StatusBadRequest, nil)
		return
	}

	// ULIDを生成
	id, err := generateULID()
	if err != nil {
		logAndSendError(w, "Failed to generate ULID", http.StatusInternalServerError, err)
		return
	}

	// 挿入用のSQLクエリを作成
	stmt, err := db.Prepare("INSERT INTO items (id, title, content, category, chapter, file, createdBy, createdByName) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logAndSendError(w, "Failed to prepare SQL statement", http.StatusInternalServerError, err)
		return
	}

	// データベースにデータを挿入
	_, err = stmt.Exec(id, data.Title, data.Content, data.Category, data.Chapter, data.File, data.CreatedBy, data.CreatedByName)
	if err != nil {
		logAndSendError(w, "Failed to execute SQL statement", http.StatusInternalServerError, err)
		return
	}

	// 成功時のレスポンスを返す
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // フロントエンドのオリジン
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	responseData := map[string]string{"id": id}
	json.NewEncoder(w).Encode(responseData)
}
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クロスオリジンリクエスト用のヘッダーを設定
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // フロントエンドのオリジン
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // クロスオリジンで許可するHTTPメソッド
		w.Header().Set("Content-Type", "application/json")

		// preflightリクエストの場合、200 OKを返して終了
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 通常のリクエストの場合、次のハンドラに処理を委譲
		next.ServeHTTP(w, r)
	})
}

// ULIDを生成する関数
func generateULID() (string, error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	ulid, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return "", err
	}
	return ulid.String(), nil
}
