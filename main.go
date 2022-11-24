package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"

	"cloud.google.com/go/pubsub"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var listChatID []int64

func main() {
	ctx := context.Background()
	go handleCommand()

	f, err := ioutil.ReadFile("default.tmpl")
	if err != nil {
		log.Print(err)
		return
	}

	bot, err := tgbotapi.NewBotAPI("5849872871:AAEtgwHnnXURk9dTttUEIRe26C-_EW6Kq60")
	if err != nil {
		log.Panic(err)
	}
	client, err := pubsub.NewClient(ctx, "my-project-1668668140794")
	if err != nil {
		// TODO: Handle error.
	}

	sub := client.Subscription("sub_one")

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		t := template.Must(template.New("").Parse(string(f)))
		in := []byte(string(m.Data))
		var raw map[string]interface{}
		if err := json.Unmarshal(in, &raw); err != nil {
			panic(err)
		}
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, raw); err != nil {
			panic(err)
		}
		fmt.Println(listChatID)
		for _, id := range listChatID {
			fmt.Println(id)
			msg := tgbotapi.NewMessage(int64(id), tpl.String())
			bot.Send(msg)
		}
		m.Ack()
	})

	if err != context.Canceled {
		// TODO: Handle error.
	}
}

func handleCommand() {
	bot, err := tgbotapi.NewBotAPI("5849872871:AAEtgwHnnXURk9dTttUEIRe26C-_EW6Kq60")
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		fmt.Println("Hello")
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			listChatID = append(listChatID, int64(update.Message.Chat.ID))
			fmt.Println(listChatID)
			msg.Text = "I understand /sayhi."
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
