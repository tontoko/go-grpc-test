Dockerfile: Dockerイメージの作成に使用するファイル。airのインストール、設定ファイルのコピー、実行コマンドなどが記述されている。

- go-grpc-test/
  - プロジェクトのルートディレクトリ
- Dockerfile
  - Dockerイメージの定義ファイル
- docker-compose.yml
  - Dockerコンテナーの定義ファイル
- book.proto
  - gRPCサービス定義ファイル
- server/
  - gRPCサーバーの実装
  - main.go
    - サーバーのエントリーポイント
- client/
  - gRPCクライアントの実装
  - main.go
    - クライアントのエントリーポイント
