package handlers

import (
	"db/database"
	"encoding/json"
	"net/http"
)

func HandleGetCategoryNames(w http.ResponseWriter) {
	rows, err := database.Db.Query("SELECT Name FROM categories")
	if err != nil {
		http.Error(w, "データベースのクエリエラー", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			http.Error(w, "行のスキャンエラー", http.StatusInternalServerError)
			return
		}
		names = append(names, name)
	}

	response, err := json.Marshal(names)
	if err != nil {
		http.Error(w, "JSONのマーシャリングエラー", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
