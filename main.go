package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/aleibovici/cryptopump/algorithms"
	"github.com/aleibovici/cryptopump/exchange"
	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/markets"
	"github.com/aleibovici/cryptopump/mysql"
	"github.com/aleibovici/cryptopump/node"
	"github.com/aleibovici/cryptopump/plotter"
	"github.com/aleibovici/cryptopump/telegram"
	"github.com/aleibovici/cryptopump/threads"
	"github.com/aleibovici/cryptopump/types"

	"github.com/jtaczanowski/go-scheduler"
	"github.com/sdcoffey/techan"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

		logger.LogEntry{
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

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
		SymbolFiat:           "",
		SymbolFiatFunds:      0,
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

	open.Run("http://localhost:" + port) /* Open URI using the OS's default browser */

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}

func (fh *myHandler) handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("X-Content-Type-Options", "nosniff") /* Add X-Content-Type-Options header */
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Add("X-Frame-Options", "DENY") /* Prevent page from being displayed in an frame */

	fh.configData = functions.GetConfigData(fh.sessionData)

	switch r.Method {
	case "GET":

		/* Determine the URI path to de taken */
		switch r.URL.Path {
		case "/":

			loadConfigDataAdditionalComponents(fh.configData, fh.sessionData) /* Load dynamic components in configData */

			functions.ExecuteTemplate(w, fh.configData, fh.sessionData) /* This is the template execution for 'index' */

		case "/sessiondata":

			var tmp []byte
			var err error

			w.Header().Set("Content-Type", "application/json")

			if tmp, err = loadSessionDataAdditionalComponents(fh.sessionData, fh.marketData, fh.configData); err != nil { /* Load dynamic components for javascript autoloader for html output */

				logger.LogEntry{
					Config:   fh.configData,
					Market:   fh.marketData,
					Session:  fh.sessionData,
					Order:    &types.Order{},
					Message:  functions.GetFunctionName() + " - " + err.Error(),
					LogLevel: "DebugLevel",
				}.Do()

				return

			}

			if _, err := w.Write(tmp); err != nil { /* Write writes the data to the connection as part of an HTTP reply. */

				logger.LogEntry{
					Config:   fh.configData,
					Market:   fh.marketData,
					Session:  fh.sessionData,
					Order:    &types.Order{},
					Message:  functions.GetFunctionName() + " - " + err.Error(),
					LogLevel: "DebugLevel",
				}.Do()

				return

			}

		}

	case "POST":

		/* Determine the URI path to de taken */
		switch r.URL.Path {
		case "/":

			/* This function reads and parse the html form */
			if err := r.ParseForm(); err != nil {

				logger.LogEntry{
					Config:   fh.configData,
					Market:   nil,
					Session:  fh.sessionData,
					Order:    &types.Order{},
					Message:  functions.GetFunctionName() + " - " + err.Error(),
					LogLevel: "DebugLevel",
				}.Do()

				return

			}

			/* This function uses a hidden field 'submitselect' in each HTML template to detect the actions triggered by users.
			HTML action must include 'document.getElementById('submitselect').value='about';this.form.submit()' */
			switch r.PostFormValue("submitselect") {
			case "new":

				/* Spawn a new process  */
				path, err := os.Executable()
				if err != nil {

					logger.LogEntry{
						Config:   fh.configData,
						Market:   nil,
						Session:  fh.sessionData,
						Order:    &types.Order{},
						Message:  functions.GetFunctionName() + " - " + err.Error(),
						LogLevel: "DebugLevel",
					}.Do()

				}

				cmd := exec.Command(path)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err = cmd.Start()
				if err != nil {

					logger.LogEntry{
						Config:   fh.configData,
						Market:   nil,
						Session:  fh.sessionData,
						Order:    &types.Order{},
						Message:  functions.GetFunctionName() + " - " + err.Error(),
						LogLevel: "DebugLevel",
					}.Do()

				}

				functions.ExecuteTemplate(w, fh.configData, fh.sessionData) /* This is the template execution for 'index' */

			case "start":

				go execution(
					fh.configData,
					fh.sessionData,
					fh.marketData)

				time.Sleep(2 * time.Second)          /* Sleep time to wait for ThreadID to start */
				http.Redirect(w, r, r.URL.Path, 301) /* Redirect to root 'index' */

			case "stop":

				threads.ExitThreadID(fh.sessionData) /* Cleanly exit ThreadID */

			case "update":

				functions.SaveConfigData(r, fh.sessionData) /* Save updated config */

				http.Redirect(w, r, r.URL.Path, 301) /* Redirect to root 'index' */

			case "buy":

				fh.sessionData.ForceBuy = true

				http.Redirect(w, r, r.URL.Path, 301) /* Redirect to root 'index' */

			case "sell":

				fh.sessionData.ForceSell = true

				http.Redirect(w, r, r.URL.Path, 301) /* Redirect to root 'index' */

			case "configTemplate":

				fh.sessionData.ConfigTemplate = functions.StrToInt(r.PostFormValue("configTemplateList")) /* Retrieve Configuration Template Key selection */

				configData := functions.LoadConfigTemplate(fh.sessionData) /* Load and populate html with Configuration Template */

				functions.ExecuteTemplate(w, configData, fh.sessionData) /* This is the template execution for 'index' */

			}
		}

	}

}

func execution(
	configData *types.Config,
	sessionData *types.Session,
	marketData *types.Market) {

	var err error

	/* Connect to Exchange */
	if err = exchange.GetClient(configData, sessionData); err != nil {

		logger.LogEntry{
			Config:   configData,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		/* Cleanly exit ThreadID */
		threads.ExitThreadID(sessionData)

	}

	/* Routine to resume operations */
	var threadIDSessionDB string

	if sessionData.ThreadID, threadIDSessionDB, err = mysql.GetThreadTransactionDistinct(sessionData); err != nil {

		logger.LogEntry{
			Config:   configData,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		/* Cleanly exit ThreadID */
		threads.ExitThreadID(sessionData)

	}

	if sessionData.ThreadID != "" && !configData.NewSession {

		/* If GetThreadTransactionDistinct return empty, create and lock thread file */
		threads.LockThreadID(sessionData.ThreadID)

		configData = functions.GetConfigData(sessionData)

		if sessionData.Symbol, err = mysql.GetOrderSymbol(sessionData); err != nil {

			logger.LogEntry{
				Config:   configData,
				Market:   nil,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "DebugLevel",
			}.Do()

			/* Cleanly exit ThreadID */
			threads.ExitThreadID(sessionData)

		}

		if sessionData.Symbol == "" {

			logger.LogEntry{
				Config:   configData,
				Market:   nil,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  "sessionData.Symbol not found",
				LogLevel: "DebugLevel",
			}.Do()

			/* Cleanly exit ThreadID */
			threads.ExitThreadID(sessionData)

		}

		/* Select the symbol coin to be used from sessionData.Symbol option */
		sessionData.SymbolFiat = sessionData.Symbol[3:7]

		logger.LogEntry{
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  "Resuming",
			LogLevel: "InfoLevel",
		}.Do()

	} else {

		/* Define Thread ID for the node */
		sessionData.ThreadID = functions.GetThreadID()

		/* Create lock for threadID */
		if !threads.LockThreadID(sessionData.ThreadID) {

			os.Exit(1)

		}

		/* Select the symbol coin to be used from Config option */
		sessionData.Symbol = configData.Symbol
		sessionData.SymbolFiat = configData.SymbolFiat

		logger.LogEntry{
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  "Initializing",
			LogLevel: "InfoLevel",
		}.Do()

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
	node.GetRole(configData, sessionData)
	scheduler.RunTaskAtInterval(
		func() { node.GetRole(configData, sessionData) },
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

	/* Update exchange latency every 5 seconds. */
	scheduler.RunTaskAtInterval(
		func() {
			functions.GetExchangeLatency(sessionData)
		},
		time.Second*5,
		time.Second*0)

	/* Check system status every 10 seconds. */
	scheduler.RunTaskAtInterval(
		func() {
			node.CheckStatus(configData, sessionData)
		},
		time.Second*10,
		time.Second*0)

	/* Retrieve available fiat funds and update database
	This is only used for retrieving balances for the first time, ans is then followed by
	the Websocket routine to retrieve realtime user data  */
	if sessionData.SymbolFiatFunds, err = exchange.GetSymbolFiatFunds(
		configData,
		sessionData); err == nil {
		_ = mysql.UpdateSession(
			configData,
			sessionData)
	}

	/* Retrieve available symbol funds
	This is only used for retrieving balances for the first time, ans is then followed by
	the Websocket routine to retrieve realtime user data  */
	if sessionData.SymbolFunds, err = exchange.GetSymbolFunds(configData, sessionData); err != nil {

		logger.LogEntry{
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	/* Retrieve exchange lot size for ticker and store in sessionData */
	exchange.GetLotSize(configData, sessionData)

	sum := 0
	for {

		/* Check start/stop times of operation */
		if configData.TimeEnforce {

			for !functions.IsInTimeRange(configData.TimeStart, configData.TimeStop) {

				logger.LogEntry{
					Config:   configData,
					Market:   marketData,
					Session:  sessionData,
					Order:    &types.Order{},
					Message:  "Sleeping",
					LogLevel: "InfoLevel",
				}.Do()

				time.Sleep(300000 * time.Millisecond)

			}

		}

		/* Update ThreadCount */
		if sessionData.ThreadCount, err = mysql.GetThreadTransactionCount(sessionData); err != nil {

			logger.LogEntry{
				Config:   configData,
				Market:   marketData,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "DebugLevel",
			}.Do()

		}

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

		wg.Wait() /* Wait for the goroutines to finish */

		logger.LogEntry{
			Config:   configData,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  "All websocket channels stopped",
			LogLevel: "DebugLevel",
		}.Do()

		sessionData.StopWs = false /* Reset goroutine channels */

		/* Reload configuration in case of WsBookTicker broken connection */
		configData = functions.GetConfigData(sessionData)

		time.Sleep(3000 * time.Millisecond)

		logger.LogEntry{
			Config:   configData,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  "Restarting",
			LogLevel: "DebugLevel",
		}.Do()

		/* repeated forever */
		sum++

	}

}

/* Load dynamic components for javascript autoloader for html output */
func loadSessionDataAdditionalComponents(
	sessionData *types.Session,
	marketData *types.Market,
	configData *types.Config) ([]byte, error) {

	type Market struct {
		Rsi3      float64 /* Relative Strength Index for 3 periods */
		Rsi7      float64 /* Relative Strength Index for 7 periods */
		Rsi14     float64 /* Relative Strength Index for 14 periods */
		MACD      float64 /* Moving average convergence divergence */
		Price     float64 /* Market Price */
		Direction int     /* Market Direction */
	}

	type Order struct {
		OrderID  string
		Quantity float64
		Quote    float64
		Price    float64
		Target   float64
	}

	type Session struct {
		ThreadID             string  /* Unique session ID for the thread */
		SellTransactionCount float64 /* Number of SELL transactions in the last 60 minutes*/
		Symbol               string  /* Symbol */
		SymbolFunds          float64 /* Available crypto funds in exchange */
		SymbolFiat           string  /* Fiat currency funds */
		SymbolFiatFunds      float64 /* Fiat currency funds */
		ProfitThreadID       float64 /* ThreadID profit */
		Profit               float64 /* Total profit */
		ThreadCount          int     /* Thread count */
		ThreadAmount         float64 /* Thread cost amount */
		Latency              int64   /* Latency between the exchange and client */
		Orders               []Order
	}

	type Update struct {
		Market  Market
		Session Session
	}

	sessiondata := Update{}

	sessiondata.Market.Rsi3 = math.Round(marketData.Rsi3*100) / 100
	sessiondata.Market.Rsi7 = math.Round(marketData.Rsi7*100) / 100
	sessiondata.Market.Rsi14 = math.Round(marketData.Rsi14*100) / 100
	sessiondata.Market.MACD = math.Round(marketData.MACD*10000) / 10000
	sessiondata.Market.Price = math.Round(marketData.Price*1000) / 1000
	sessiondata.Market.Direction = marketData.Direction

	sessiondata.Session.Latency = sessionData.Latency /* Latency between the exchange and client */
	sessiondata.Session.ThreadID = sessionData.ThreadID
	sessiondata.Session.SellTransactionCount = sessionData.SellTransactionCount
	sessiondata.Session.Symbol = sessionData.Symbol[0:3]
	sessiondata.Session.SymbolFunds = math.Round((sessionData.SymbolFunds)*100000000) / 100000000
	sessiondata.Session.SymbolFiat = sessionData.SymbolFiat
	sessiondata.Session.SymbolFiatFunds = math.Round(sessionData.SymbolFiatFunds*100) / 100

	if profit, err := mysql.GetProfit(sessionData); err == nil {
		sessiondata.Session.Profit = math.Round(profit*100) / 100
	}
	if profitThreadID, err := mysql.GetProfitByThreadID(sessionData); err == nil {
		sessiondata.Session.ProfitThreadID = math.Round(profitThreadID*100) / 100
	}
	if threadCount, err := mysql.GetThreadCount(sessionData); err == nil {
		sessiondata.Session.ThreadCount = threadCount
	}
	if threadAmount, err := mysql.GetThreadAmount(sessionData); err == nil {
		sessiondata.Session.ThreadAmount = math.Round(threadAmount*100) / 100
	}

	if orders, err := mysql.GetThreadTransactionByThreadID(sessionData); err == nil {

		for _, key := range orders {

			tmp := Order{}
			tmp.OrderID = strconv.Itoa(key.OrderID)
			tmp.Quantity = key.ExecutedQuantity
			tmp.Quote = math.Round(key.CumulativeQuoteQuantity*100) / 100
			tmp.Price = math.Round(key.Price*10000) / 10000
			tmp.Target = math.Round((tmp.Price*(1+configData.ProfitMin))*1000) / 1000

			sessiondata.Session.Orders = append(sessiondata.Session.Orders, tmp)
		}

	}

	return json.Marshal(sessiondata)

}

/* Load dynamic components into configData for html output */
func loadConfigDataAdditionalComponents(
	configData *types.Config,
	sessionData *types.Session) {

	configData.HTMLSnippet = plotter.Plot(sessionData)

}
