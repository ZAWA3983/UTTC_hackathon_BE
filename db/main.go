// main.go

package main

import (
	"db/cors"
	"db/database"
	"db/handlers"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

func init() {
	database.Init()
}

func main() {
	// CORSミドルウェアを適用
	http.Handle("/api/addItem", cors.CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 通常のリクエストの処理
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			// preflightリクエストの場合、200 OKを返して終了
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodPost:
			// POSTリクエストの処理
			handlers.HandleAddItem(w, r)
		default:
			// サポートされていないメソッドの場合、405 Method Not Allowedを返す
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	http.Handle("/api/searchItems", cors.CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 通常のリクエストの処理
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			// preflightリクエストの場合、200 OKを返して終了
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodPost:
			// POSTリクエストの処理
			handlers.HandleSearchItems(w, r)
		default:
			// サポートされていないメソッドの場合、405 Method Not Allowedを返す
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	http.Handle("/api/myItems", cors.CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 通常のリクエストの処理
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			// preflightリクエストの場合、200 OKを返して終了
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodPost:
			// POSTリクエストの処理
			handlers.HandleSearchMyItems(w, r)
		default:
			// サポートされていないメソッドの場合、405 Method Not Allowedを返す
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})))

	http.Handle("/api/updateItem", cors.CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodPut:
			handlers.HandleUpdateItems(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))

	http.Handle("/api/deleteItem", cors.CORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodDelete:
			handlers.HandleDeleteItem(w, r)
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
