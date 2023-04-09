# Chat GPT APIリクエスト

テキストチャットと画像生成のリクエストをするgolangのコードです。

それぞれベースはchat gpt自身に作って貰いました。
ただそのままでは動かなかったので、jsonの内容などをAPIリファレンスの内容に改めたり、環境変数でAPIキーを指定するようにしてみたりしています。

## 環境変数

環境変数```OPENAI_API_KEY```に取得したAPIキーを設定しておいてください。

## 実行

```shell
go run src/gogptchat/gogptchat.go "こんにちは"
go run src/gogptimg/gogptimg.go "馬の絵を描いてください"
```
画像生成の方は、生成された画像のURLが2つ提示されます。
