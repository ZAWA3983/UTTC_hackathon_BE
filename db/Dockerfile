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

# 環境変数 PORT からポート番号を取得
ARG PORT=8080
ENV PORT=$ {PORT:-8080} 

# ポートの公開
EXPOSE $PORT
