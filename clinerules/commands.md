airのインストール: go install github.com/cosmtrek/air@latest

- make docker-build
  - Dockerイメージをビルドする
- make docker-up
  - Dockerコンテナーを起動する
- make docker-down
  - Dockerコンテナーを停止する
- make proto-gen
  - Protocol BuffersからGoコードを生成する
- go mod init go-grpc-test
  - Goモジュールを初期化する