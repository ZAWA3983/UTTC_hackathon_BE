package handlers

import (
	"db/database"
	"db/model"
	"encoding/json"
	"net/http"
	"time"
)

// HandleSearchMyItems はPOSTリクエストを処理する関数
func HandleSearchMyItems(w http.ResponseWriter, r *http.Request) {
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
	rows, err := database.Db.Query(sqlQuery, params...)
	if err != nil {
		logAndSendError(w, "Failed to execute SQL query", http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	// 結果をスライスにマップ
	var items []model.Item
	for rows.Next() {
		var item model.Item
		var createdAtStr string // DATETIME 型のデータを文字列として読み込む
		var updatedAtStr string
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.Category,
			&item.Chapter,
			&item.File,
			&item.FileType, // fileType を読み込む
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
