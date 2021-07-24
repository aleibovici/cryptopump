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
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/jtaczanowski/go-scheduler"
	"github.com/sdcoffey/techan"
	"github.com/skratchdot/open-golang/open"
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

	open.Run("http://localhost:" + port) /* Open URI using the OS's default browser */

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

func (fh *myHandler) handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("X-Content-Type-Options", "nosniff") /* Add X-Content-Type-Options header */
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Add("X-Frame-Options", "DENY") /* Prevent page from being displayed in an iframe */

	fh.configData = functions.GetConfigData(fh.sessionData)

	switch r.Method {
	case "GET":

		/* Determine the URI path to de taken */
		switch r.URL.Path {
		case "/":

			functions.ExecuteTemplate(w, fh.configData, fh.sessionData) /* This is the template execution for 'index' */

		case "/sessiondata":

			/* Load dynamic components for javascript autoloader for html output */

			w.Header().Set("Content-Type", "application/json")

			tmp, _ := loadSessionDataAdditionalComponents(fh.sessionData, fh.marketData, fh.configData)

			if _, err := w.Write(tmp); err != nil {

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
			case "new":

				/* Spawn a new CryptoPump process  */
				path, err := os.Executable()
				if err != nil {
					log.Println(err)
				}

				cmd := exec.Command(path)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err = cmd.Start()
				if err != nil {
					log.Println(err)
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

/* Load dynamic components for javascript autoloader for html output */
func loadSessionDataAdditionalComponents(
	sessionData *types.Session,
	marketData *types.Market,
	configData *types.Config) ([]byte, error) {

	type Market struct {
		Rsi3        string /* Relative Strength Index for 3 periods */
		Rsi7        string /* Relative Strength Index for 7 periods */
		Rsi14       string /* Relative Strength Index for 14 periods */
		MACD        string /* Moving average convergence divergence */
		Price       string /* Market Price */
		Direction   string /* Market Direction */
		HtmlSnippet template.HTML
	}

	type Order struct {
		OrderID string
		Quote   string
		Price   string
	}

	type Session struct {
		ThreadID             string /* Unique session ID for the thread */
		SellTransactionCount string /* Number of SELL transactions in the last 60 minutes*/
		Symbol_fiat          string /* Fiat currency funds */
		Symbol_fiat_funds    string /* Fiat currency funds */
		ProfitThreadID       string /* ThreadID profit */
		Profit               string /* Total profit */
		ThreadCount          string /* Thread count */
		ThreadAmount         string /* Thread cost amount */
		Orders               []Order
	}

	type Update struct {
		Market  Market
		Session Session
	}

	sessiondata := Update{}

	sessiondata.Market.Rsi3 = functions.Float64ToStr(marketData.Rsi3, 2)
	sessiondata.Market.Rsi7 = functions.Float64ToStr(marketData.Rsi7, 2)
	sessiondata.Market.Rsi14 = functions.Float64ToStr(marketData.Rsi14, 2)
	sessiondata.Market.MACD = functions.Float64ToStr(marketData.MACD, 4)
	sessiondata.Market.Price = functions.Float64ToStr(marketData.Price, 3)
	sessiondata.Market.Direction = strconv.Itoa(marketData.Direction)

	sessiondata.Session.ThreadID = sessionData.ThreadID
	sessiondata.Session.SellTransactionCount = functions.Float64ToStr(sessionData.SellTransactionCount, 0)
	sessiondata.Session.Symbol_fiat = sessionData.Symbol_fiat
	sessiondata.Session.Symbol_fiat_funds = functions.Float64ToStr(sessionData.Symbol_fiat_funds, 2)

	if profit, err := mysql.GetProfit(sessionData); err == nil {
		sessiondata.Session.Profit = functions.Float64ToStr(profit, 2)
	}
	if profitThreadID, err := mysql.GetProfitByThreadID(sessionData); err == nil {
		sessiondata.Session.ProfitThreadID = functions.Float64ToStr(profitThreadID, 2)
	}
	if threadCount, err := mysql.GetThreadCount(sessionData); err == nil {
		sessiondata.Session.ThreadCount = strconv.Itoa(threadCount)
	}
	if threadAmount, err := mysql.GetThreadAmount(sessionData); err == nil {
		sessiondata.Session.ThreadAmount = functions.Float64ToStr(threadAmount, 2)
	}

	if orders, err := mysql.GetThreadTransactionByThreadID(sessionData); err == nil {

		for _, key := range orders {

			tmp := Order{}
			tmp.OrderID = strconv.Itoa(key.OrderID)
			tmp.Quote = functions.Float64ToStr(key.CummulativeQuoteQuantity, 2)
			tmp.Price = functions.Float64ToStr(key.Price, 3)

			sessiondata.Session.Orders = append(sessiondata.Session.Orders, tmp)
		}

	}

	sessiondata.Market.HtmlSnippet = plotter.Plot(sessionData)

	return json.Marshal(sessiondata)

}
