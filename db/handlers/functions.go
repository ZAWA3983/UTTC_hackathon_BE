package handlers

import (
	"github.com/oklog/ulid"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// ULIDを生成する関数
func generateULID() (string, error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	ulid, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return "", err
	}
	return ulid.String(), nil
}

// 指定されたエラーメッセージをログに記録してHTTPエラーレスポンスを返すユーティリティ関数
func logAndSendError(w http.ResponseWriter, message string, status int, err error) {
	log.Printf("Error: %v\n", err)
	http.Error(w, message, status)
}
