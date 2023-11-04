# ベースとなるDockerイメージを選択
FROM golang:1.18

# コンテナ内の作業ディレクトリを設定
WORKDIR /app

# ソースコードをコピー
COPY main.go .
COPY go.mod .
COPY go.sum .

# Goモジュールをインストール
RUN go mod download

# Goアプリケーションをビルド
RUN go build -o app main.go

# アプリケーションの起動コマンドを設定
CMD ["./app"]

# ポート8080を公開
EXPOSE 8080
