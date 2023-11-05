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

	http.Handle("/api/searchItems", CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 通常のリクエストの処理
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			// preflightリクエストの場合、200 OKを返して終了
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodPost:
			// POSTリクエストの処理
			handleSearchItems(w, r)
		default:
			// サポートされていないメソッドの場合、405 Method Not Allowedを返す
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	http.Handle("/api/myItems", CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 通常のリクエストの処理
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			// preflightリクエストの場合、200 OKを返して終了
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodPost:
			// POSTリクエストの処理
			handleSearchItems(w, r)
		default:
			// サポートされていないメソッドの場合、405 Method Not Allowedを返す
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})))

	http.Handle("/api/updateItem", CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodPut:
			handleUpdateItem(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	http.Handle("/api/deleteItem", CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodDelete:
			handleDeleteItem(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server listening on :" + port)
	http.ListenAndServe(":"+port, nil)
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
	w.Header().Set("Access-Control-Allow-Origin", "https://uttc-hackathon-fe.vercel.app") // フロントエンドのオリジン
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	responseData := map[string]string{"id": id}
	json.NewEncoder(w).Encode(responseData)
}

// handleSearchItems はPOSTリクエストを処理する関数
func handleSearchItems(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドがPOSTでない場合はエラーを返す
	if r.Method != http.MethodPost {
		logAndSendError(w, "Only POST requests are allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	// リクエストボディからデータをデコード
	var queryData struct {
		SearchTerm string `json:"searchTerm"`
		Category   string `json:"category"`
		Chapter    string `json:"chapter"`
		SortOption string `json:"sortOption"`
	}
	if err := json.NewDecoder(r.Body).Decode(&queryData); err != nil {
		logAndSendError(w, "Failed to decode request body", http.StatusBadRequest, err)
		return
	}

	// SQLクエリを構築
	// ソートオプションに応じて適切なORDER BY句を追加
	sortSQL := ""
	switch queryData.SortOption {
	case "createdAt":
		sortSQL = "ORDER BY createdAt DESC"
	case "-createdAt":
		sortSQL = "ORDER BY createdAt"
	case "updatedAt":
		sortSQL = "ORDER BY updatedAt DESC"
	case "-updatedAt":
		sortSQL = "ORDER BY updatedAt"
	}

	// パラメータ化されたSQLクエリを構築
	sqlQuery := "SELECT * FROM items WHERE title LIKE ?"
	params := []interface{}{"%" + queryData.SearchTerm + "%"}

	// カテゴリと章の選択肢が空でない場合、それらをクエリに追加
	if queryData.Category != "" {
		sqlQuery += " AND category = ?"
		params = append(params, queryData.Category)
	}
	if queryData.Chapter != "" {
		sqlQuery += " AND chapter = ?"
		params = append(params, queryData.Chapter)
	}

	sqlQuery += " " + sortSQL // ソートオプションを適用

	// SQLクエリを実行
	rows, err := db.Query(sqlQuery, params...)
	if err != nil {
		logAndSendError(w, "Failed to execute SQL query", http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	// 結果をスライスにマップ
	var items []Item
	for rows.Next() {
		var item Item
		var createdAtStr string // DATETIME 型のデータを文字列として読み込む
		var updatedAtStr string
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.Category,
			&item.Chapter,
			&item.File,
			&item.CreatedBy,
			&item.CreatedByName,
			&createdAtStr, // 文字列として読み込む
			&updatedAtStr,
		)
		if err != nil {
			logAndSendError(w, "Failed to scan row", http.StatusInternalServerError, err)
			return
		}

		// createdAt と updatedAt の文字列を time.Time に変換
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			logAndSendError(w, "Failed to parse createdAt", http.StatusInternalServerError, err)
			return
		}
		updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			logAndSendError(w, "Failed to parse updatedAt", http.StatusInternalServerError, err)
			return
		}

		item.CreatedAt = createdAt
		item.UpdatedAt = updatedAt
		items = append(items, item)
	}

	// 検索結果をJSONレスポンスとして返す
	w.Header().Set("Access-Control-Allow-Origin", "https://uttc-hackathon-fe.vercel.app") // フロントエンドのオリジン
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クロスオリジンリクエスト用のヘッダーを設定
		w.Header().Set("Access-Control-Allow-Origin", "https://uttc-hackathon-fe.vercel.app") // フロントエンドのオリジン
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST,PUT,DELETE, OPTIONS") // クロスオリジンで許可するHTTPメソッド
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

// MyItems はPOSTリクエストを処理する関数
func MyItems(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドがPOSTでない場合はエラーを返す
	if r.Method != http.MethodPost {
		logAndSendError(w, "Only POST requests are allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	// リクエストボディからデータをデコード
	var queryData struct {
		SearchTerm string `json:"searchTerm"`
		Category   string `json:"category"`
		Chapter    string `json:"chapter"`
		SortOption string `json:"sortOption"`
		UserEmail  string `json:"userEmail"` // ユーザーのメールアドレスを追加
	}
	if err := json.NewDecoder(r.Body).Decode(&queryData); err != nil {
		logAndSendError(w, "Failed to decode request body", http.StatusBadRequest, err)
		return
	}

	// SQLクエリを構築
	// ソートオプションに応じて適切なORDER BY句を追加
	sortSQL := ""
	switch queryData.SortOption {
	case "createdAt":
		sortSQL = "ORDER BY createdAt DESC"
	case "-createdAt":
		sortSQL = "ORDER BY createdAt"
	case "updatedAt":
		sortSQL = "ORDER BY updatedAt DESC"
	case "-updatedAt":
		sortSQL = "ORDER BY updatedAt"
	}

	// パラメータ化されたSQLクエリを構築
	sqlQuery := "SELECT * FROM items WHERE title LIKE ?"
	params := []interface{}{"%" + queryData.SearchTerm + "%"}

	// ユーザーのメールアドレスを条件に追加（CreatedBy との一致）
	sqlQuery += " AND createdBy = ?"
	params = append(params, queryData.UserEmail)

	// カテゴリと章の選択肢が空でない場合、それらをクエリに追加
	if queryData.Category != "" {
		sqlQuery += " AND category = ?"
		params = append(params, queryData.Category)
	}
	if queryData.Chapter != "" {
		sqlQuery += " AND chapter = ?"
		params = append(params, queryData.Chapter)
	}

	sqlQuery += " " + sortSQL // ソートオプションを適用

	// SQLクエリを実行
	rows, err := db.Query(sqlQuery, params...)
	if err != nil {
		logAndSendError(w, "Failed to execute SQL query", http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	// 結果をスライスにマップ
	var items []Item
	for rows.Next() {
		var item Item
		var createdAtStr string // DATETIME 型のデータを文字列として読み込む
		var updatedAtStr string
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.Category,
			&item.Chapter,
			&item.File,
			&item.CreatedBy,
			&item.CreatedByName,
			&createdAtStr, // 文字列として読み込む
			&updatedAtStr,
		)
		if err != nil {
			logAndSendError(w, "Failed to scan row", http.StatusInternalServerError, err)
			return
		}

		// createdAt と updatedAt の文字列を time.Time に変換
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			logAndSendError(w, "Failed to parse createdAt", http.StatusInternalServerError, err)
			return
		}
		updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			logAndSendError(w, "Failed to parse updatedAt", http.StatusInternalServerError, err)
			return
		}

		item.CreatedAt = createdAt
		item.UpdatedAt = updatedAt
		items = append(items, item)
	}

	// 検索結果をJSONレスポンスとして返す
	w.Header().Set("Access-Control-Allow-Origin", "https://uttc-hackathon-fe.vercel.app") // フロントエンドのオリジン
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

// handleUpdateItem はPUTリクエストを処理する関数
func handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		logAndSendError(w, "Only PUT requests are allowed", http.StatusMethodNotAllowed, nil)
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

	// アイテムを更新するSQLクエリを作成
	stmt, err := db.Prepare("UPDATE items SET title = ?, content = ?, category = ?, chapter = ?, file = ?, createdByName = ?, updatedAt = NOW() WHERE id = ?")
	if err != nil {
		logAndSendError(w, "Failed to prepare SQL statement", http.StatusInternalServerError, err)
		return
	}

	// データベースのアイテムを更新
	_, err = stmt.Exec(data.Title, data.Content, data.Category, data.Chapter, data.File, data.CreatedByName, data.ID)
	if err != nil {
		logAndSendError(w, "Failed to execute SQL statement", http.StatusInternalServerError, err)
		return
	}

	// 成功時のレスポンスを返す
	w.Header().Set("Access-Control-Allow-Origin", "https://uttc-hackathon-fe.vercel.app") // フロントエンドのオリジン
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseData := map[string]string{"message": "更新が成功しました"}
	json.NewEncoder(w).Encode(responseData)
}

func handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		logAndSendError(w, "Only DELETE requests are allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	// リクエストボディから削除対象のアイテムIDを取得
	var data struct {
		ItemIds []string `json:"itemIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logAndSendError(w, "Failed to decode request body", http.StatusBadRequest, err)
		return
	}

	// アイテムを削除するSQLクエリを実行
	for _, itemId := range data.ItemIds {
		_, err := db.Exec("DELETE FROM items WHERE id = ?", itemId)
		if err != nil {
			logAndSendError(w, "Failed to execute SQL statement", http.StatusInternalServerError, err)
			return
		}
	}

	// 削除が成功した場合のレスポンスを返す
	w.Header().Set("Access-Control-Allow-Origin", "https://uttc-hackathon-fe.vercel.app")
	w.WriteHeader(http.StatusOK)
	responseData := map[string]string{"message": "削除が成功しました"}
	json.NewEncoder(w).Encode(responseData)
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
