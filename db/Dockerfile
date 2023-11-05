# ベースとなるDockerイメージを選択
FROM golang:1.21

# コンテナ内の作業ディレクトリを設定
WORKDIR /app

# ソースコードをコピー
COPY main.go .
COPY go.mod .
COPY go.sum .

# Goアプリケーションをビルド
RUN go build -o app main.go

# アプリケーションの起動コマンドを設定
CMD ["./app"]

# ポート8080の公開
EXPOSE 8080
