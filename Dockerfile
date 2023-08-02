# Dockerfileは、Dockerイメージの作成手順を指定するテキストドキュメント。
# ベースのOS、追加するファイル、実行するコマンド、開くポートなど、全体的な設定を定義する。
# `docker build` コマンドで使用される。

# Goの公式Dockerイメージ
# https://hub.docker.com/_/golang

# ------------------------------------------------------------------------------

# リリース用のビルドを行う、コンテナイメージ作成ステージ
# リリース用のコンテナイメージには含みたくない秘匿情報を含んだファイルや環境変数を利用可能
FROM golang:1.20.6-bullseye as deploy-builder

WORKDIR /app

# アプリの依存関係をコピー
COPY go.mod go.sum ./
# 依存関係のダウンロード
RUN go mod download

# アプリのソースをバンドルする
COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# ------------------------------------------------------------------------------

# マネージドサービス上で動かすことを想定した、リリース用のコンテナイメージ作成ステージ
# docker build -t Riscait/go_todo_app:${DOCKER_TAG} --target deploy ./
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

# アプリを実行
CMD ["./app"]

# ------------------------------------------------------------------------------

# ローカル開発環境（ホットリロード可能）で利用する、コンテナイメージ作成ステージ
FROM golang:1.20.6 as dev

WORKDIR /app

# Goでホットリロード開発を実現するOSSをインストール
RUN go install github.com/cosmtrek/air@latest

# airを実行
CMD [ "air" ]
