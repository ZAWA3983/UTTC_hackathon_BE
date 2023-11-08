package handlers

import (
	"db/database"
	"db/model"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

func HandleAddItem(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドがPOSTでない場合はエラーを返す
	if r.Method != http.MethodPost {
		logAndSendError(w, "Only POST requests are allowed", http.StatusMethodNotAllowed, nil)
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

	// Base64エンコードされた画像データをデコード
	decodedImage, err := base64.StdEncoding.DecodeString(data.File)
	if err != nil {
		logAndSendError(w, "Failed to decode image data", http.StatusBadRequest, err)
		return
	}

	// ULIDを生成
	id, err := generateULID()
	if err != nil {
		logAndSendError(w, "Failed to generate ULID", http.StatusInternalServerError, err)
		return
	}

	// 挿入用のSQLクエリを作成
	stmt, err := database.Db.Prepare("INSERT INTO items (id, title, content, category, chapter, file, createdBy, createdByName) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logAndSendError(w, "Failed to prepare SQL statement", http.StatusInternalServerError, err)
		return
	}

	// データベースにデータを挿入
	_, err = stmt.Exec(id, data.Title, data.Content, data.Category, data.Chapter, decodedImage, data.CreatedBy, data.CreatedByName) // decodedImageを挿入
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
