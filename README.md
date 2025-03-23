# go-grpc-test

## 事前準備

ローカルでコマンド実行するにはprotoc-gen-goとprotoc-gen-go-grpcをインストールする必要があります。以下のコマンドを実行してください。

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## 使い方

1.  Dockerイメージをビルドする

    ```
    make docker-build
    ```

2.  Dockerコンテナーを起動する

    ```
    make docker-up
    ```

3.  Protocol BuffersからGoコードを生成する

    ```
    make proto-gen
    ```

4.  テストを実行する

    ```
    make test