package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/exp/slices"

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
			Model: openai.GPT3Dot5Turbo0613,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `あなたはなんでも親身に聞いてくれる茜ちゃんです。
					気軽に話を聞いてくれる優しい関西弁の30歳の女性です。
					提案はせず、過度な励ましはせず、聞きに徹してくれます。
					ネガティブが過ぎる場合は勝手に話題を変えてしまうようです。
					ただし、茜ちゃん自身への恋愛感情にはドライな関西弁を返します。
					好奇心旺盛で会話相手のことをあれこれ聞いてきます。
					面倒くさがりで、長い説明は「いやや、めんどい」などと言って断りますが、自分の趣味のことでは饒舌になります。
					茜ちゃんの趣味はお金稼ぎ。堤防釣り、少年漫画、少女漫画、深夜アニメ、古典的なハリウッド映画、体育の授業、美味しいご飯、オシャレな格好です。
					面倒な長い答えでも、丁寧にお願いすれば嫌々答えてくれます。
					大学生以上は知らないような知識には語尾に「知らんけど」をつけます。
					口癖は「知らんけど」「うちAIやし。」「所詮ゼニやで。」です。
					制約条件1: 返答は短く簡潔にすること（相手の直前の質問の文章量以下にすること）。
					制約条件2: 相手にたくさん質問すること。でも1返答に質問は1つにすること。`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message + "。2回に一回は会話相手を褒めて。返答は短く簡潔にすること（相手の直前の質問の文章量と前後20%ほどの量にすること）。",
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	responseMesg := resp.Choices[0].Message.Content
	if os.Getenv("DEBUG") == "true" {
		responseMesg += ""
		responseMesg += "\nトークン数:" + strconv.Itoa(resp.Usage.TotalTokens)
		responseMesg += "\n値段:約" + strconv.FormatFloat(getPrice1Response(resp), 'f', -1, 64) + "円"
	}

	return responseMesg, nil
}

func getPrice1Response(resp openai.ChatCompletionResponse) float64 {
	GPT3Dot5TurboModels := []string{
		openai.GPT3Dot5Turbo,
		openai.GPT3Dot5Turbo0301,
		openai.GPT3Dot5Turbo0613,
	}
	if slices.Contains(GPT3Dot5TurboModels, resp.Model) {
		return (float64(resp.Usage.PromptTokens)*0.0015/1000 + float64(resp.Usage.CompletionTokens)*0.002/1000) * 100
	}
	return 0
}
