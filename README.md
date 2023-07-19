# aichatbot
AIが返答を考えるLINEチャットボットです。

# 開発環境での実行方法
LINEのトークン, シークレットは[LINE Developersコンソール](https://developers.line.biz/console/)と[LINE Official Account Manager](https://manager.line.biz/)を確認する。
OpenAI APIのキーは設定時に発行したものを利用すること。

ターミナル1
```
export PORT=8080
export LINE_CHANNEL_TOKEN=LINE API トークン
export LINE_CHANNEL_SECRET=LINE API シークレット
export OPENAI_API_KEY=OpenAI APIキー
go run api/main.go
```

ターミナル2
```
ngrok http 8080
```

ブラウザで[LINE Developersコンソール](https://developers.line.biz/console/)のWebhookURLにngrokで払い出したURLに`/callback`をつけて設定する
