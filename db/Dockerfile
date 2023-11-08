# ベースとなるDockerイメージを選択 (Goのバージョンを指定)
FROM golang:1.21

# コンテナ内の作業ディレクトリを設定
WORKDIR /app

# ソースコードをコピー
COPY . .

# Goアプリケーションをビルド
RUN go build -o app main.go

# ポート番号を環境変数から取得
ARG PORT=8080
ENV PORT=$PORT

# ポート番号によってエクスポーズするポートを設定
EXPOSE $PORT

# アプリケーションの起動コマンドを設定
CMD ["./app"]
