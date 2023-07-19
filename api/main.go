package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"

	openai "github.com/sashabaranov/go-openai"
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
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_TOKEN"),
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
				// 返答文章の取得
				replyMessage, merr := getResMessage(message.Text)
				if merr != nil {
					fmt.Println(merr.Error())
					return
				}
				// 返答
				_, rerr := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(replyMessage),
				).Do()
				if rerr != nil {
					fmt.Println(rerr.Error())
				}
			}
		}
	}
}

func getResMessage(message string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
