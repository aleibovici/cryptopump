package main

import (
	"cryptopump/algorithms"
	"cryptopump/exchange"
	"cryptopump/functions"
	"cryptopump/markets"
	"cryptopump/mysql"
	"cryptopump/node"
	"cryptopump/plotter"
	"cryptopump/telegram"
	"cryptopump/threads"
	"cryptopump/types"
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jtaczanowski/go-scheduler"
	"github.com/sdcoffey/techan"
	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type myHandler struct {
	sessionData *types.Session
	marketData  *types.Market
	configData  *types.Config
}

func init() {

	viper.SetConfigName("config")   /* Set the file name of the configurations file */
	viper.AddConfigPath("./config") /* Set the path to look for the configurations file */
	viper.SetConfigType("yml")      /* */
	viper.AutomaticEnv()            /* Enable VIPER to read Environment Variables */

	if err := viper.ReadInConfig(); err != nil {

		functions.Logger(
			nil,
			nil,
			nil,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

	}

	viper.WatchConfig()

}

func main() {

	sessionData := &types.Session{
		ThreadID:             "",
		ThreadIDSession:      "",
		ThreadCount:          0,
		SellTransactionCount: 0,
		Symbol:               "",
		Symbol_fiat:          "",
		Symbol_fiat_funds:    0,
		LastBuyTransactTime:  time.Time{},
		LastSellCanceledTime: time.Time{},
		ConfigTemplate:       0,
		ForceBuy:             false,
		ForceSell:            false,
		ListenKey:            "",
		MasterNode:           false,
		TgBotAPI:             &tgbotapi.BotAPI{},
		Db:                   &sql.DB{},
		Clients:              types.Client{},
		KlineData:            []types.KlineData{},
		StopWs:               false,
		Busy:                 false,
		MinQuantity:          0,
		MaxQuantity:          0,
		StepSize:             0,
	}

	marketData := &types.Market{
		Rsi3:                      0,
		Rsi7:                      0,
		Rsi14:                     0,
		MACD:                      0,
		Price:                     0,
		PriceChangeStatsHighPrice: 0,
		PriceChangeStatsLowPrice:  0,
		Direction:                 0,
		TimeStamp:                 time.Time{},
		Series:                    &techan.TimeSeries{},
	}

	configData := &types.Config{}

	/* Initialize DB connection */
	sessionData.Db = mysql.DBInit()

	myHandler := &myHandler{
		sessionData: sessionData,
		marketData:  marketData,
		configData:  configData}

	port := functions.GetPort() /* Determine port for HTTP service. */

	http.HandleFunc("/", myHandler.handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Printf("Listening on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

/* Load dynamic components into configData interface for html output */
func loadConfigDataAdditionalComponents(
	configData *types.Config,
	sessionData *types.Session,
	marketData *types.Market) {

	configData.HtmlSnippet = plotter.Plot(sessionData)
	configData.Orders, _ = mysql.GetThreadTransactionByThreadID(sessionData)   /* Update open orders for html output */
	configData.FiatFunds = math.Round(sessionData.Symbol_fiat_funds*100) / 100 /* Store fiat currency funds for html output */
	configData.Profit, _ = mysql.GetProfit(sessionData)                        /* Update total profit for html output */
	configData.ProfitThreadID, _ = mysql.GetProfitByThreadID(sessionData)      /* Update thread profit for html output */
	configData.SellTransactionCount = sessionData.SellTransactionCount         /* Store Number of SELL transactions in the last 60 minutes for html output */
	configData.ThreadCount, _ = mysql.GetThreadCount(sessionData)              /* Store thread count for html output */
	configData.ThreadAmount, _ = mysql.GetThreadAmount(sessionData)            /* Store thread cost amount for html output */
	configData.MarketDataMACD = math.Floor(marketData.MACD*100) / 100          /* Store  for html output */
	configData.MarketDataRsi14 = math.Floor(marketData.Rsi14*100) / 100        /* Store RSI14 for html output */
	configData.MarketDataRsi7 = math.Floor(marketData.Rsi7*100) / 100          /* Store RSI7 for html output */
	configData.MarketDataRsi3 = math.Floor(marketData.Rsi3*100) / 100          /* Store RSI3 for html output */

}

func (fh *myHandler) handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("X-Content-Type-Options", "nosniff") // Add X-Content-Type-Options header
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Add("X-Frame-Options", "DENY") // Prevent page from being displayed in an iframe

	fh.configData = functions.GetConfigData(fh.sessionData)

	switch r.Method {
	case "GET":

		/* Determine the URI path to de taken */
		switch r.URL.Path {
		case "/":

			/* Load dynamic components into configData interface for html output */
			loadConfigDataAdditionalComponents(fh.configData, fh.sessionData, fh.marketData)

			/* This is the template execution for 'index' */
			functions.ExecuteTemplate(w, "index.html", fh.configData, fh.sessionData)

		}

	case "POST":

		/* Determine the URI path to de taken */
		switch r.URL.Path {

		case "/":

			/* This function reads and parse the html form */
			if err := r.ParseForm(); err != nil {

				functions.Logger(
					fh.configData,
					nil,
					fh.sessionData,
					log.DebugLevel,
					0,
					0,
					0,
					0,
					functions.GetFunctionName()+" - "+err.Error())

				return

			}

			/* This function uses a hidden field 'submitselect' in each HTML template to detect the actions triggered by users.
			HTML action must include 'document.getElementById('submitselect').value='about';this.form.submit()' */
			switch r.PostFormValue("submitselect") {
			case "start":

				go execution(
					fh.configData,
					fh.sessionData,
					fh.marketData)

				/* Load dynamic components into configData interface for html output */
				loadConfigDataAdditionalComponents(fh.configData, fh.sessionData, fh.marketData)

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", fh.configData, fh.sessionData)

			case "stop":

				/* Cleanly exit ThreadID */
				threads.ExitThreadID(fh.sessionData)

			case "update":

				/* Save updated config */
				functions.SaveConfigData(r, fh.sessionData)

				/* Load dynamic components into configData interface for html output */
				loadConfigDataAdditionalComponents(fh.configData, fh.sessionData, fh.marketData)

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", fh.configData, fh.sessionData)

			case "buy":

				fh.sessionData.ForceBuy = true

				/* Load dynamic components into configData interface for html output */
				loadConfigDataAdditionalComponents(fh.configData, fh.sessionData, fh.marketData)

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", fh.configData, fh.sessionData)

			case "sell":

				fh.sessionData.ForceSell = true

				/* Load dynamic components into configData interface for html output */
				loadConfigDataAdditionalComponents(fh.configData, fh.sessionData, fh.marketData)

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", fh.configData, fh.sessionData)

			case "configTemplate":

				/* Retrieve Configuration Template Key selection */
				fh.sessionData.ConfigTemplate = functions.StrToInt(r.PostFormValue("configTemplateList"))

				/* Load and populate html with Configuration Template */
				configData := functions.LoadConfigTemplate(fh.sessionData)

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", configData, fh.sessionData)

			}
		}

	}

}

func execution(
	configData *types.Config,
	sessionData *types.Session,
	marketData *types.Market) {

	var err error

	/* Define Exchange to be used */
	exchange.GetClient(configData, sessionData)

	/* Routine to resume operations */
	var threadIDSessionDB string
	sessionData.ThreadID, threadIDSessionDB, _ = mysql.GetThreadTransactionDistinct(sessionData)

	if sessionData.ThreadID != "" && configData.NewSession == "false" {

		configData = functions.GetConfigData(sessionData)

		sessionData.Symbol, _ = mysql.GetOrderSymbol(sessionData)

		if sessionData.Symbol == "" {

			functions.Logger(
				configData,
				nil,
				sessionData,
				log.InfoLevel,
				0,
				0,
				0,
				0,
				"sessionData.Symbol not found")

			/* Cleanly exit ThreadID */
			threads.ExitThreadID(sessionData)

		}

		/* Select the symbol coin to be used from sessionData.Symbol option */
		sessionData.Symbol_fiat = sessionData.Symbol[3:7]

		functions.Logger(
			configData,
			marketData,
			sessionData,
			log.InfoLevel,
			0,
			0,
			0,
			0,
			"Resuming")

	} else {

		/* Define Thread ID for the node */
		sessionData.ThreadID = functions.GetThreadID()

		/* Create lock for threadID */
		if !functions.LockThreadID(sessionData.ThreadID) {

			os.Exit(1)

		}

		/* Select the symbol coin to be used from Config option */
		sessionData.Symbol = configData.Symbol.(string)
		sessionData.Symbol_fiat = configData.Symbol_fiat.(string)

		functions.Logger(
			configData,
			marketData,
			sessionData,
			log.InfoLevel,
			0,
			0,
			0,
			0,
			"Initializing")

	}

	/* Print threadID to debug for easy identification of session */
	fmt.Printf("ThreadID:  %s", sessionData.ThreadID)

	/* Synchronize time with Binance every 5 minutes */
	_ = exchange.NewSetServerTimeService(configData, sessionData)
	scheduler.RunTaskAtInterval(
		func() { _ = exchange.NewSetServerTimeService(configData, sessionData) },
		time.Second*300,
		time.Second*0)

	/* Retrieve config data every 10 seconds. */
	scheduler.RunTaskAtInterval(
		func() { configData = functions.GetConfigData(sessionData) },
		time.Second*10,
		time.Second*0)

	/* run function UpdatePendingOrders() every 180 seconds */
	rand.Seed(time.Now().UnixNano())
	scheduler.RunTaskAtInterval(
		func() { algorithms.UpdatePendingOrders(configData, sessionData) },
		time.Second*180,
		time.Second*time.Duration(rand.Intn(180-1+1)+1),
	)

	/* Retrieve initial node role and then every 60 seconds */
	node.GetRole(sessionData)
	scheduler.RunTaskAtInterval(
		func() { node.GetRole(sessionData) },
		time.Second*60,
		time.Second*0)

	/* Keep user stream service alive every 60 seconds */
	scheduler.RunTaskAtInterval(
		func() { _ = exchange.KeepAliveUserStreamServiceListenKey(configData, sessionData) },
		time.Second*60,
		time.Second*0)

	/* Update Number of Sale Transactions per hour every 3 minutes.
	The same function is executed after each sale, and when initiating cycle. */
	scheduler.RunTaskAtInterval(
		func() {
			sessionData.SellTransactionCount, err = mysql.GetOrderTransactionCount(sessionData, "SELL")
		},
		time.Second*180,
		time.Second*0)

	/* Retrieve available fiat funds and update database
	This is only used for retrieving balances for the first time, ans is then followed by
	the Websocket routine to retrieve realtime user data  */
	if sessionData.Symbol_fiat_funds, _ = exchange.GetSymbolFunds(
		configData,
		sessionData); err == nil {
		_ = mysql.UpdateSession(
			configData,
			sessionData)
	}

	/* Retrieve exchange lot size for ticker and store in sessionData */
	exchange.GetLotSize(configData, sessionData)

	sum := 0
	for {

		/* Check start/stop times of operation */
		if configData.Time_enforce.(string) == "true" {

			for !functions.IsInTimeRange(configData.Time_start.(string), configData.Time_stop.(string)) {

				functions.Logger(
					configData,
					marketData,
					sessionData,
					log.InfoLevel,
					0,
					0,
					0,
					0,
					"Sleeping")

				time.Sleep(300000 * time.Millisecond)

			}

		}

		/* Update ThreadCount */
		sessionData.ThreadCount, _ = mysql.GetThreadTransactionCount(sessionData)

		/* Update Number of Sale Transactions per hour */
		sessionData.SellTransactionCount, err = mysql.GetOrderTransactionCount(sessionData, "SELL")

		/* This routine is executed when no transaction cycle has initiated (ThreadCount = 0) */
		if sessionData.ThreadCount == 0 {

			/* Define new Thread ID Session */
			sessionData.ThreadIDSession = functions.GetThreadID()

			/* Save new session to Session table. */
			if err := mysql.SaveSession(
				configData,
				sessionData); err != nil {

				/* Update existing session on Session table */
				if err := mysql.UpdateSession(
					configData,
					sessionData); err != nil {

					/* Cleanly exit ThreadID */
					threads.ExitThreadID(sessionData)

				}

			}

		} else {

			/* Retrieve existing Thread ID Session if first time */
			if threadIDSessionDB != "" {

				sessionData.ThreadIDSession = threadIDSessionDB
				threadIDSessionDB = ""

				/* Save new session to Session table then update if fail */
				if err := mysql.SaveSession(
					configData,
					sessionData); err != nil {

					/* Update existing session on Session table */
					if err := mysql.UpdateSession(
						configData,
						sessionData); err != nil {

						/* Cleanly exit ThreadID */
						threads.ExitThreadID(sessionData)

					}

				}

			}

		}

		/* Conditional used in case this is the first run in the cycle go get past market data */
		if marketData.PriceChangeStatsHighPrice == 0 {

			markets.LoadKlineDataPast(
				configData,
				marketData,
				sessionData)

		}

		wg := &sync.WaitGroup{} /* WaitGroup to stop inside Channels */
		wg.Add(3)               /* WaitGroup to stop inside Channels */

		/* Start Telegram bot if Master Node and store in sessionData.TgBotAPI */
		go telegram.CheckUpdates(
			configData,
			sessionData,
			wg)

		/* Websocket routine to retrieve realtime candle data */
		go algorithms.WsKline(
			configData,
			marketData,
			sessionData,
			wg)

		/* Websocket routine to retrieve realtime user data */
		go algorithms.WsUserDataServe(
			configData,
			sessionData,
			wg)

		/* Websocket routine to retrieve realtime ticker prices */
		go algorithms.WsBookTicker(
			configData,
			marketData,
			sessionData,
			wg)

		wg.Wait()                  /* Wait for the goroutines to finish */
		sessionData.StopWs = false /* Reset goroutine channels */

		/* Reload configuration in case of WsBookTicker broken connection */
		configData = functions.GetConfigData(sessionData)

		time.Sleep(3000 * time.Millisecond)

		/* repeated forever */
		sum++

	}

}
