package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

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

func roller(amount int, diff int, expl string) string {

	var res string
	rand.Seed(time.Now().UnixNano())
	var iter int = 0
	var ran int
	res = "Amount = " + strconv.Itoa(amount) + "\n" + "diff " + strconv.Itoa(diff)
	res = res + "\n" + "expl = " + expl + "\n" + "numbers:"
	for i := 0; i < amount; i++ {
		ran = rand.Intn(10)
		if ran == 0 {
			ran = 10
		}
		res = res + " " + strconv.Itoa(ran)
		if ran >= diff {
			iter++
			if (ran == 10) && (expl == "y") {
				ran = rand.Intn(11)
				if ran == 0 {
					ran = 10
				}
				res = res + " (exploded)" + strconv.Itoa(ran)
				if ran >= diff {
					iter++
				}
			}
		}
		if ran == 1 {
			iter--
		}
	}
	res = res + "\n result:" + strconv.Itoa(iter)
	return res
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
	if param[0] <= 0 || param[1] <= 0 || (ret[2] != "y" && ret[2] != "n") {
		return "Bad params"
	}
	return res + roller(param[0], param[1], ret[2])

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
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upd, _ := bot.GetUpdatesChan(ucfg)
	// читаем обновления из канала
	for x != 1 {
		x = 0
		select {
		case update := <-upd:
			// Пользователь, который написал боту
			if !update.Message.IsCommand() {
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = helper(update.Message.CommandArguments())
			case "roll":
				msg.Text = standard_roll(update.Message.CommandArguments())
			case "exit":
				msg.Text = "Good Night " + update.Message.Chat.UserName
				if update.Message.Chat.UserName == "HedonisticAI" {
					x = 1
				}
			default:
				msg.Text = "I don't know that command type /help for info"
			}
			bot.Send(msg)
		}

	}
}
