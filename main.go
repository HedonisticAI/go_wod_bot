package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func helper(Text string) string {
	switch Text {
	case "roll":
		return "roll is usually run with three params\n first amount of dice\n second - difficulty \n y/n for explosion"
	case "exit":
		return "type fo end conversation"
	default:
		return "I understand /roll and /exit.For commands info type /help with roll or exit"
	}

}

// 0 transforms to 10 because
func dicerand(Text string) string {
	var ret = strings.Split(Text, " ")
	var res string
	rand.Seed(time.Now().UnixNano())
	var amount, _ = strconv.Atoi(ret[0])
	var diff, _ = strconv.Atoi(ret[1])
	var iter int = 0
	var ran int
	res = "Amount = " + ret[0] + "\n" + "diff " + ret[1]
	res = res + "\n" + "expl = " + ret[2] + "\n" + "numbers:"
	for i := 0; i < amount; i++ {
		ran = rand.Intn(10)
		if ran == 0 {
			ran = 10
		}
		res = res + " " + strconv.Itoa(ran)
		if ran >= diff {
			iter++
			if (ran == 10) && (ret[2] == "y") {
				ran = rand.Intn(11)
				res = res + " " + strconv.Itoa(ran)
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
				msg.Text = dicerand(update.Message.CommandArguments())
			case "exit":
				msg.Text = "Good Night " + update.Message.Chat.UserName
				x = 1
			default:
				msg.Text = "I don't know that command type /help for info"
			}
			bot.Send(msg)
		}

	}
}
