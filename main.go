package main

import (
	// Библиотеки, нужные программе
	"github.com/porty/go-osx-screenshot"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Горутина, запускающая заствку
func screensaver() {
	exec.Command("/System/Library/Frameworks/ScreenSaver.framework/Versions/A/Resources/ScreenSaverEngine.app/Contents/MacOS/ScreenSaverEngine").Output()
}

func main() {
	//	if exec.Command("ping -c4 8.8.8.8") != nil {
	//		log.Printf("Ping OK")
	//	}
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("234464051:AAH50Nox97nDpcJCWH9ekl_y3yn9wcfeY8E")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	// читаем обновления из канала
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		// Созадаем сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		img := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, "")
		// Ping-Pong
		if update.Message.Text == "/ping" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
			// Пингуем сайты
		} else if update.Message.Text[0:4] == "ping" {
			site := strings.TrimPrefix(update.Message.Text, "ping ")
			out, err := exec.Command("ping", "-c4", site).Output()
			if err != nil {
				panic(err)
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, string(out))
			// Запускаем заставку из горутины
		} else if update.Message.Text == "/off" {
			go screensaver()
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
			// Отрпака скриншота
		} else if update.Message.Text == "/screenshot" {

			// Задаём имя файла с датой
			ttime := time.Now()
			tt := ttime.Format("2006.01.02.15.04")
			filename := "Screenshots/screenshot." + tt + ".png"
			// Делаем скриншот
			var format screenshot.SaveFormat
			if strings.HasSuffix(strings.ToLower(filename), ".png") {
				format = screenshot.FormatPng
			}
			err := screenshot.SaveScreenshotToFile(filename, format)
			if err != nil {
				panic(err)
			}
			// отправляем скриншот
			img = tgbotapi.NewDocumentUpload(update.Message.Chat.ID, filename)
			img.ReplyToMessageID = update.Message.MessageID
			bot.Send(img)
			// В остальных случаях
		} else if update.Message != nil {
			// Ответить на сообщение его копией
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		}
		// и отправляем его
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

	}
}
