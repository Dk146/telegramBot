package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"text/template"

	"cloud.google.com/go/pubsub"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx := context.Background()

	botToken := os.Getenv("BOT_TOKEN")
	templatePath := os.Getenv("TEMPLATE")
	pubsubProjectId := os.Getenv("PUBSUB_PROJECT_ID")
	pubsubSubId := os.Getenv("PUBUB_SUB_ID")
	chatId := os.Getenv("CHAT_ID")

	f, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	client, err := pubsub.NewClient(ctx, pubsubProjectId)
	if err != nil {
		log.Panic(err)
	}
	sub := client.Subscription(pubsubSubId)

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		t := template.Must(template.New("").Parse(string(f)))
		in := []byte(string(m.Data))
		fmt.Println(string(m.Data))

		var raw map[string]interface{}
		if err := json.Unmarshal(in, &raw); err != nil {
			fmt.Println(err)
		}
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, raw); err != nil {
			fmt.Println(err)
		}

		int_chat_id, _ := strconv.Atoi(chatId)
		msg := tgbotapi.NewMessage(int64(int_chat_id), tpl.String())
		bot.Send(msg)
		m.Ack()
	})

	if err != context.Canceled {
		log.Panic(err)
	}
}
