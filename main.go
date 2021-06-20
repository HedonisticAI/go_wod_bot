package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type Config struct {
	TelegramBotToken string
}

func helper(Text string) string {
	switch Text {
	case "roll":
		return "roll is usually run with three params\n first amount of dice\n second - difficulty \n y/n for explosion"
	case "exit":
		return "type to of the bot(for dev)"
	default:
		return "I understand /roll and /exit.For commands info type /help with roll or exit"
	}

}

// 0 transforms to 10 because when you roll standard d10 you in contains 0...9 and you count 0 as 10
func standard_roll(Text string) string {
	if Text == "" {
		return "No params"
	}
	var ret = strings.Split(Text, " ")
	var res string
	var param [2]int
	param[0], _ = strconv.Atoi(ret[0])
	param[1], _ = strconv.Atoi(ret[1])
	if !checkstd(param[0], param[1], ret[2]) {
		return "Bad params"
	}
	return res + roller(param[0], param[1], ret[2])

}

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi, i am roller and assistant"))
}

func main() {
	file, _ := os.Open("cfg.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	var x int8 = 0
	// инициализируем канал, куда будут прилетать обновления от API
	//var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	//ucfg.Timeout = 60
	//upd, _ := bot.GetUpdatesChan(ucfg)
	upd := bot.ListenForWebhook("/" + bot.Token)
	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	// читаем обновления из канала
	for x != 1 {
		x = 0
		select {
		case update := <-upd:
			// Пользователь, который написал боту
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = "default"
			if update.Message.IsCommand() {
				msg.Text = commandhandler(update)
				if update.Message.Command() == "exit" {
					x = 1
				}
				msg.ReplyToMessageID = update.Message.MessageID
			}

			bot.Send(msg)
		default:
			continue
		}

	}
}
