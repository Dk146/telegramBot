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

	bot_token := os.Getenv("BOT_TOKEN")
	template_path := os.Getenv("TEMPLATE")
	pubsub_project_id := os.Getenv("PUBSUB_PROJECT_ID")
	pubsub_sub_id := os.Getenv("PUBUB_SUB_ID")
	chat_id := os.Getenv("CHAT_ID")

	f, err := ioutil.ReadFile(template_path)
	if err != nil {
		log.Panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(bot_token)
	if err != nil {
		log.Panic(err)
	}
	client, err := pubsub.NewClient(ctx, pubsub_project_id)
	if err != nil {
		log.Panic(err)
	}
	sub := client.Subscription(pubsub_sub_id)

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

		int_chat_id, _ := strconv.Atoi(chat_id)
		msg := tgbotapi.NewMessage(int64(int_chat_id), tpl.String())
		bot.Send(msg)
		m.Ack()
	})

	if err != context.Canceled {
		log.Panic(err)
	}
}
