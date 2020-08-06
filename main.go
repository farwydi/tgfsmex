package main

import (
	"awesomeProject3/xfsm"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("668993758:AAEDSNWEkiBE55jQtvoiuAYmgpDQzYx-W_c")
	if err != nil {
		log.Fatalln(err)
	}

	transition := []xfsm.Transition{
		// W
		{Name: "/start", Source: "W", Destination: "W"},
		{Name: ">begin", Source: "W", Destination: "Q0"},
		// R
		{Name: "/start", Source: "R", Destination: "R"},
		{Name: ">reset", Source: "R", Destination: "W"},
	}

	state := []xfsm.State{
		{
			Name: "W",
			Callback: func(message *tgbotapi.Message) error {
				msg := tgbotapi.NewMessage(
					message.Chat.ID,
					"Welcome page",
				)
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Begin", ">begin"),
					),
				)
				_, err := bot.Send(msg)
				return err
			},
		},
		{
			Name: "R",
			Callback: func(message *tgbotapi.Message) error {
				msg := tgbotapi.NewMessage(
					message.Chat.ID,
					"Result page",
				)
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Reset", ">reset"),
					),
				)
				_, err := bot.Send(msg)
				return err
			},
		},
	}

	countQ := 5
	for i := 0; i < countQ; i++ {
		stateName := fmt.Sprintf("Q%d", i)
		nextStateName := fmt.Sprintf("Q%d", i+1)
		if i+1 == countQ {
			nextStateName = "R"
		}
		transition = append(transition, xfsm.Transition{Name: "/start", Source: stateName, Destination: stateName})
		transition = append(transition, xfsm.Transition{Name: ">a1", Source: stateName, Destination: nextStateName})
		transition = append(transition, xfsm.Transition{Name: ">a2", Source: stateName, Destination: nextStateName})

		state = append(state, xfsm.State{
			Name: stateName,
			Callback: func(message *tgbotapi.Message) error {
				msg := tgbotapi.NewMessage(
					message.Chat.ID,
					stateName+" page",
				)
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("A1", ">a1"),
						tgbotapi.NewInlineKeyboardButtonData("A2", ">a2"),
					),
				)
				_, err := bot.Send(msg)
				return err
			},
		})
	}

	machine := xfsm.NewFSM(
		"W",
		transition,
		state,
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	log.Println("Ready")

	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
		fmt.Println(machine.Current())
		if update.CallbackQuery != nil {
			err := machine.Event(update.CallbackQuery.Data, update.CallbackQuery.Message)
			if err != nil {
				fmt.Println(err)
			}
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				err := machine.Event(update.Message.Text, update.Message)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				log.Println("cmd or callback only")
			}
		} else {
			log.Println("???")
		}
	}
}
