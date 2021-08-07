package telegram

import (
	"cryptopump/functions"
	"cryptopump/logger"
	"cryptopump/mysql"
	"cryptopump/threads"
	"cryptopump/types"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

/* Establish connectivity to Telegram */
func connect(
	configData *types.Config,
	sessionData *types.Session) (tgBotAPI *tgbotapi.BotAPI) {

	var err error

	if tgBotAPI, err = tgbotapi.NewBotAPI(configData.TgBotApikey); err != nil {

		logger.LogEntry{
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	tgBotAPI.Debug = false

	return tgBotAPI

}

/* Send message to Telegram */
func send(
	message tgbotapi.MessageConfig,
	sessionData *types.Session) {

	if _, err := sessionData.TgBotAPI.Send(message); err != nil {

		logger.LogEntry{
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()
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
	if configData.TgBotApikey == "" {

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

		logger.LogEntry{
			Config:   configData,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

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

		case "/report":

			var profit float64
			var threadCount int
			var err error

			if profit, err = mysql.GetProfit(sessionData); err != nil {
				return
			}

			if threadCount, err = mysql.GetThreadCount(sessionData); err != nil {
				return
			}

			tmp := "\f" + "Funds: " + sessionData.SymbolFiat + " " + functions.Float64ToStr(sessionData.SymbolFiatFunds, 2) + "\n" +
				"Profit: " + functions.Float64ToStr(profit, 2) + "\n" +
				"Thread Count: " + strconv.Itoa(threadCount) + "\n" +
				"Master: " + sessionData.ThreadID

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, tmp)
			msg.ReplyToMessageID = update.Message.MessageID
			send(msg, sessionData)

		}

	}

}
