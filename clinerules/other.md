.air.toml: airの設定ファイル。ファイルの監視対象、再起動時のコマンド、環境変数などを設定する。
docker-compose.yml: Docker Composeの設定ファイル。アプリケーションのビルド、ポート、再起動ポリシー、ボリュームなどを設定する。


- その他
  - Docker環境で動作させることを前提とする
  - Go Modulesを使用して依存関係を管理する
  - docker-compose.ymlに`restart: always`を設定
    - コンテナーが予期せず停止した場合でも自動的に再起動されるように設定
    - `docker-compose up --build`コマンドでコンテナーを再起動する必要がある