package telegram

import (
	"math"
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
}

// Connect to connect to Telegram
type Connect struct{}

// Send a message via Telegram
func (message Message) Send(sessionData *types.Session) {

	msg := tgbotapi.NewMessage(sessionData.TgBotAPIChatID, message.Text)

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

// Do establish connectivity to Telegram
func (Connect) Do(
	configData *types.Config,
	sessionData *types.Session) {

	var err error

	if sessionData.TgBotAPI, err = tgbotapi.NewBotAPI(configData.TgBotApikey); err != nil {

		logger.LogEntry{
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	sessionData.TgBotAPI.Debug = false

}

// CheckUpdates Check for Telegram bot updates
func CheckUpdates(
	configData *types.Config,
	sessionData *types.Session,
	wg *sync.WaitGroup) {

	var err error
	var updates tgbotapi.UpdatesChannel

	/* Exit if no API key found */
	if configData.TgBotApikey == "" {

		return

	}

	/* Sleep until Master Node is True */
	for !sessionData.MasterNode {

		time.Sleep(30000 * time.Millisecond)

	}

	/* Establish connectivity to Telegram server */
	Connect{}.Do(configData, sessionData)

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

		/* Store Telegram ChatID to allow the system to send direct messages to Telegram server */
		sessionData.TgBotAPIChatID = update.Message.Chat.ID

		switch update.Message.Text {
		case "/sell":

			Message{
				Text:             "\f" + "Selling @ " + sessionData.ThreadID,
				ReplyToMessageID: update.Message.MessageID,
			}.Send(sessionData)

			sessionData.ForceSell = true

		case "/buy":

			Message{
				Text:             "\f" + "Buying @ " + sessionData.ThreadID,
				ReplyToMessageID: update.Message.MessageID,
			}.Send(sessionData)

			sessionData.ForceBuy = true

		case "/report":

			var profit float64
			var profitPct float64
			var roi float64
			var threadCount int
			var status string
			var err error

			if profit, profitPct, err = mysql.GetProfit(sessionData); err != nil {
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

			/* Calculate total Return on Investment */
			roi = (profit) / (math.Round(sessionData.Global.ThreadAmount*100) / 100) * 100

			Message{
				Text: "\f" + "Available Funds: " + sessionData.SymbolFiat + " " + functions.Float64ToStr(sessionData.SymbolFiatFunds, 2) + "\n" +
					"Deployed Funds: " + sessionData.SymbolFiat + " " + functions.Float64ToStr((math.Round(sessionData.Global.ThreadAmount*100)/100), 2) + "\n" +
					"Net Profit (Profit - Deployed Diff): " + functions.Float64ToStr(profit, 2) + " " + functions.Float64ToStr(profitPct, 2) + "%" + "\n" +
					"ROI " + functions.Float64ToStr(roi, 2) + "%" + "\n" +
					"Thread Count: " + strconv.Itoa(threadCount) + "\n" +
					"Status: " + status + "\n" +
					"Master: " + sessionData.ThreadID,
				ReplyToMessageID: update.Message.MessageID,
			}.Send(sessionData)

		}

	}

}
