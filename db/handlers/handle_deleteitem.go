package handlers

import (
	"db/database"
	"encoding/json"
	"net/http"
)

func HandleDeleteItem(w http.ResponseWriter, r *http.Request) {
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
		_, err := database.Db.Exec("DELETE FROM items WHERE id = ?", itemId)
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
