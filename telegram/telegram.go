package telegram

import (
	"strconv"
	"sync"
	"time"

	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/mysql"
	"github.com/aleibovici/cryptopump/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Message defines the message structure to send via Telegram
type Message struct {
	Text             string
	ReplyToMessageID int
	ChatID           int64
}

// Send message via Telegram
func (message Message) Send(sessionData *types.Session) {

	msg := tgbotapi.NewMessage(message.ChatID, message.Text)

	if _, err := sessionData.TgBotAPI.Send(msg); err != nil {

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
			// threads.ExitThreadID(sessionData)

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
			var status string
			var err error

			if profit, err = mysql.GetProfit(sessionData); err != nil {
				return
			}

			if threadCount, err = mysql.GetThreadCount(sessionData); err != nil {
				return
			}

			if threadID, err := mysql.GetSessionStatus(sessionData); err == nil {

				if threadID != "" {
					status = "\f" + "System Fault @ " + threadID
				} else {
					status = "\f" + "System nominal"
				}

			}

			Message{
				Text: "\f" + "Funds: " + sessionData.SymbolFiat + " " + functions.Float64ToStr(sessionData.SymbolFiatFunds, 2) + "\n" +
					"Total Profit: " + functions.Float64ToStr(profit, 2) + "\n" +
					"Thread Count: " + strconv.Itoa(threadCount) + "\n" +
					"Status: " + status + "\n" +
					"Master: " + sessionData.ThreadID,
				ChatID:           update.Message.Chat.ID,
				ReplyToMessageID: update.Message.MessageID,
			}.Send(sessionData)

			/* Store Telegram chat ID to allow the system to send updates to user */
			sessionData.TgBotAPIChatID = update.Message.Chat.ID

		}

	}

}
