package handlers

import (
	"db/database"
	"db/model"
	"encoding/json"
	"net/http"
)

// HandleUpdateItems はPUTリクエストを処理する関数
func HandleUpdateItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		logAndSendError(w, "Only PUT requests are allowed", http.StatusMethodNotAllowed, nil)
		return
	}

	// リクエストボディからデータをデコード
	var data model.Item
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
	stmt, err := database.Db.Prepare("UPDATE items SET title = ?, content = ?, category = ?, chapter = ?, file = ?, createdByName = ?, updatedAt = NOW() WHERE id = ?")
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
