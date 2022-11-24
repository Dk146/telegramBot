package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"text/template"

	"cloud.google.com/go/pubsub"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx := context.Background()

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
		msg := tgbotapi.NewMessage(1280014428, tpl.String())
		bot.Send(msg)
		m.Ack()
	})

	if err != context.Canceled {
		// TODO: Handle error.
	}
}
