package main

import (
	"math/rand"
	"strconv"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// checker for standard roll
func checkstd(amount int, diff int, expl string) bool {
	if amount <= 0 || diff <= 0 || (expl != "y" && expl != "n") {
		return false
	}
	if amount > 1000 || diff > 10 {
		return false
	}
	return true
}

//0 transforms to 10 because on standard dice we have 0...9
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
				if ran == 1 {
					iter--
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

func commandhandler(update tgbotapi.Update) string {
	switch update.Message.Command() {
	case "help":
		return helper(update.Message.CommandArguments())
	case "roll":
		return standard_roll(update.Message.CommandArguments())
	case "exit":
		return "Good Night " + update.Message.Chat.UserName
	}
	return "Unknown Command"
}
