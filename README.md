# aichatbot
AIが返答を考えるLINEチャットボットです。

# 開発環境での実行方法
トークン, シークレットは[LINE Developersコンソール](https://developers.line.biz/console/)と[LINE Official Account Manager](https://manager.line.biz/)を確認する。

ターミナル1
```
export PORT=8080
export CHANNEL_TOKEN=トークン
export CHANNEL_SECRET=シークレット
go run api/main.go
```

ターミナル2
```
ngrok http 8080
```

ブラウザで[LINE Developersコンソール](https://developers.line.biz/console/)のWebhookURLにngrokで払い出したURLに`/callback`をつけて設定する
