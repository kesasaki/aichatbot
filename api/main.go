package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	engine := gin.Default()
	engine.GET("/", getTop)
	engine.POST("/callback", postCallback)
	engine.Run(":" + os.Getenv("PORT"))
}

func getTop(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}

func postCallback(c *gin.Context) {
	// bot作成
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// リクエスト処理
	events, berr := bot.ParseRequest(c.Request)
	if berr != nil {
		fmt.Println(berr.Error())
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				_, rerr := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(getResMessage(message.Text)),
				).Do()
				if rerr != nil {
					fmt.Println(rerr.Error())
				}
			}
		}
	}
}

func getResMessage(message string) string {
	return "あなたは" + message + "と言いました。"
}
