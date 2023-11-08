package cors

import "net/http"

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
