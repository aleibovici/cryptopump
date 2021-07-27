package telegram

import (
	"cryptopump/functions"
	"cryptopump/threads"
	"cryptopump/types"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

/* Establish connectivy to Telegram server */
func connect(
	configData *types.Config,
	sessionData *types.Session) (tgBotAPI *tgbotapi.BotAPI) {

	var err error

	if tgBotAPI, err = tgbotapi.NewBotAPI(configData.TgBotApikey.(string)); err != nil {

		functions.Logger(
			configData,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

	}

	tgBotAPI.Debug = false

	return tgBotAPI

}

/* Send message to Telegram */
func send(
	message tgbotapi.MessageConfig,
	sessionData *types.Session) {

	if _, err := sessionData.TgBotAPI.Send(message); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

	}

}

// CheckUpdates Check for Telegram bot updates
func CheckUpdates(
	configData *types.Config,
	sessionData *types.Session,
	wg *sync.WaitGroup) {

	var err error
	var msg tgbotapi.MessageConfig
	var updates tgbotapi.UpdatesChannel

	/* Exit if no API key found */
	if configData.TgBotApikey.(string) == "" {

		return

	}

	/* Sleep until Master Node is True */
	for !sessionData.MasterNode {

		time.Sleep(30000 * time.Millisecond)

	}

	/* Start Telegram bot and store in sessionData.TgBotAPI */
	sessionData.TgBotAPI = connect(
		configData,
		sessionData)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	if updates, err = sessionData.TgBotAPI.GetUpdatesChan(u); err != nil {

		functions.Logger(
			configData,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

	}

	for update := range updates {

		/* ignore any non-Message Updates */
		if update.Message == nil {

			continue

		}

		switch update.Message.Text {
		case "/stop":

			tmp := "Stopping " + sessionData.ThreadID
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, tmp)
			msg.ReplyToMessageID = update.Message.MessageID
			send(msg, sessionData)

			/* Cleanly exit ThreadID */
			threads.ExitThreadID(sessionData)

		case "/sell":

			tmp := "Selling @ " + sessionData.ThreadID
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, tmp)
			msg.ReplyToMessageID = update.Message.MessageID
			send(msg, sessionData)

			sessionData.ForceSell = true

		case "/buy":

			tmp := "Buying @ " + sessionData.ThreadID
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, tmp)
			msg.ReplyToMessageID = update.Message.MessageID
			send(msg, sessionData)

			sessionData.ForceBuy = true

		case "/funds":

			tmp := sessionData.SymbolFiat + " " + functions.Float64ToStr(sessionData.SymbolFiatFunds, 2)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, tmp)
			msg.ReplyToMessageID = update.Message.MessageID
			send(msg, sessionData)

		case "/master":

			tmp := "Master " + sessionData.ThreadID
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, tmp)
			msg.ReplyToMessageID = update.Message.MessageID
			send(msg, sessionData)

		}

	}

}
