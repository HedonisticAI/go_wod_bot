package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func dicerand(Text string) string {
	var ret = strings.Split(Text, " ")
	var res string
	rand.Seed(time.Now().UnixNano())
	return res[0]
}

func main() {
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("1840560184:AAHjCrmX8sPW-yj9mld5Qze8Mqtkn1mWdIE")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	var x int8 = 0
	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upd, _ := bot.GetUpdatesChan(ucfg)
	// читаем обновления из канала
	for x != 1 {
		select {
		case update := <-upd:
			// Пользователь, который написал боту
			UserName := update.Message.From.UserName

			// ID чата/диалога.
			// Может быть идентификатором как чата с пользователем
			// (тогда он равен UserID) так и публичного чата/канала
			ChatID := update.Message.Chat.ID

			// Текст сообщения
			Text := update.Message.Text
			if Text == "exit" {
				break
			}
			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			reply := dicerand(Text)
			// Созадаем сообщение
			msg := tgbotapi.NewMessage(ChatID, reply)
			// и отправляем его
			bot.Send(msg)
		}

	}
}
