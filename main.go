package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/aleibovici/cryptopump/algorithms"
	"github.com/aleibovici/cryptopump/exchange"
	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/loader"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/markets"
	"github.com/aleibovici/cryptopump/mysql"
	"github.com/aleibovici/cryptopump/nodes"
	"github.com/aleibovici/cryptopump/plotter"
	"github.com/aleibovici/cryptopump/telegram"
	"github.com/aleibovici/cryptopump/threads"
	"github.com/aleibovici/cryptopump/types"
	"github.com/jtaczanowski/go-scheduler"
	"github.com/paulbellamy/ratecounter"
	"github.com/sdcoffey/techan"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type myHandler struct {
	sessionData *types.Session
	marketData  *types.Market
	configData  *types.Config
	viperData   *types.ViperData
}

func main() {

	viperData := &types.ViperData{ /* Viper Configuration */
		V1: viper.New(), /* Session configurations file */
		V2: viper.New(), /* Global configurations file */
	}

	viperData.V1.SetConfigType("yml")      /* Set the type of the configurations file */
	viperData.V1.AddConfigPath("./config") /* Set the path to look for the configurations file */
	viperData.V1.SetConfigName("config")   /* Set the file name of the configurations file */
	if err := viperData.V1.ReadInConfig(); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}
	viperData.V1.WatchConfig()

	viperData.V2.SetConfigType("yml")           /* Set the type of the configurations file */
	viperData.V2.AddConfigPath("./config")      /* Set the path to look for the configurations file */
	viperData.V2.SetConfigName("config_global") /* Set the file name of the configurations file */
	if err := viperData.V2.ReadInConfig(); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}
	viperData.V2.WatchConfig()

	sessionData := &types.Session{
		ThreadID:                "",
		ThreadIDSession:         "",
		ThreadCount:             0,
		SellTransactionCount:    0,
		Symbol:                  "",
		SymbolFunds:             0,
		SymbolFiat:              "",
		SymbolFiatFunds:         0,
		LastBuyTransactTime:     time.Time{},
		LastSellCanceledTime:    time.Time{},
		LastWsKlineTime:         time.Time{},
		LastWsBookTickerTime:    time.Time{},
		LastWsUserDataServeTime: time.Time{},
		ConfigTemplate:          0,
		ForceBuy:                false,
		ForceSell:               false,
		ForceSellOrderID:        0,
		ListenKey:               "",
		MasterNode:              false,
		TgBotAPI:                &tgbotapi.BotAPI{},
		TgBotAPIChatID:          0,
		Db:                      &sql.DB{},
		Clients:                 types.Client{},
		KlineData:               []types.KlineData{},
		StopWs:                  false,
		Busy:                    false,
		MinQuantity:             0,
		MaxQuantity:             0,
		StepSize:                0,
		Latency:                 0,
		Status:                  false,
		RateCounter:             ratecounter.NewRateCounter(5 * time.Second),
		BuyDecisionTreeResult:   "",
		SellDecisionTreeResult:  "",
		QuantityOffsetFlag:      false,
		DiffTotal:               0,
		Global:                  &types.Global{},
		Admin:                   false,
		Port:                    "",
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
		Ma7:                       0,
		Ma14:                      0,
	}

	configData := &types.Config{}

	sessionData.Db = mysql.DBInit() /* Initialize DB connection */

	myHandler := &myHandler{
		sessionData: sessionData,
		marketData:  marketData,
		configData:  configData,
		viperData:   viperData,
	}

	sessionData.Port = functions.GetPort() /* Determine port for HTTP service. */

	logger.LogEntry{ /* Log Entry */
		Config:   configData,
		Market:   marketData,
		Session:  sessionData,
		Order:    &types.Order{},
		Message:  "Listening on port " + sessionData.Port,
		LogLevel: "InfoLevel",
	}.Do()

	http.HandleFunc("/", myHandler.handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	open.Run("http://localhost:" + sessionData.Port) /* Open URI using the OS's default browser */

	http.ListenAndServe(fmt.Sprintf(":%s", sessionData.Port), nil) /* Start HTTP service. */

}

func (fh *myHandler) handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")                                        /* Set the Content-Type header */
	w.Header().Set("X-Content-Type-Options", "nosniff")                                /* Add X-Content-Type-Options header */
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains") /* Add Strict-Transport-Security header */
	w.Header().Add("X-Frame-Options", "DENY")                                          /* Add X-Frame-Options header */

	fh.configData = functions.GetConfigData(fh.viperData, fh.sessionData) /* Get configuration data */

	switch r.Method {
	case "GET":

		/* Determine the URI path to de taken */
		switch r.URL.Path {
		case "/":

			fh.configData.HTMLSnippet = plotter.Data{}.Plot(fh.sessionData) /* Load dynamic components in configData */
			functions.ExecuteTemplate(w, fh.configData, fh.sessionData)     /* This is the template execution for 'index' */

		case "/sessiondata":

			var tmp []byte
			var err error

			w.Header().Set("Content-Type", "application/json") /* Set the Content-Type header */

			if tmp, err = loader.LoadSessionDataAdditionalComponents(fh.sessionData, fh.marketData, fh.configData); err != nil { /* Load dynamic components for javascript autoloader for html output */

				logger.LogEntry{ /* Log Entry */
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

				logger.LogEntry{ /* Log Entry */
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

		switch r.URL.Path { /* Determine the URI path to de taken */
		case "/":

			/* This function reads and parse the html form */
			if err := r.ParseForm(); err != nil {

				logger.LogEntry{ /* Log Entry */
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
			case "adminEnter":

				fh.sessionData.Admin = true                                 /* Set admin flag */
				functions.ExecuteTemplate(w, fh.configData, fh.sessionData) /* This is the template execution for 'admin' */

			case "adminExit":

				fh.sessionData.Admin = false                                    /* Unset admin flag */
				functions.SaveConfigGlobalData(fh.viperData, r, fh.sessionData) /* Save global data */
				functions.GetConfigData(fh.viperData, fh.sessionData)           /* Get Config Data */
				functions.ExecuteTemplate(w, fh.configData, fh.sessionData)     /* This is the template execution for 'index' */

			case "new":

				var path string /* Path to the executable */
				var err error

				/* Spawn a new process  */
				if path, err = os.Executable(); err != nil { /* Get the path of the executable */

					logger.LogEntry{ /* Log Entry */
						Config:   fh.configData,
						Market:   nil,
						Session:  fh.sessionData,
						Order:    &types.Order{},
						Message:  functions.GetFunctionName() + " - " + err.Error(),
						LogLevel: "DebugLevel",
					}.Do()

				}

				cmd := exec.Command(path) /* Spawn a new process */
				cmd.Stdout = os.Stdout    /* Redirect stdout to os.Stdout */
				cmd.Stderr = os.Stderr    /* Redirect stderr to os.Stderr */

				if err = cmd.Start(); err != nil { /* Start the new process */

					logger.LogEntry{ /* Log error */
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

				go execution(fh.viperData, fh.configData, fh.sessionData, fh.marketData) /* Start the execution process */
				time.Sleep(2 * time.Second)                                              /* Sleep time to wait for ThreadID to start */
				http.Redirect(w, r, fmt.Sprintf(r.URL.Path), 301)                        /* Redirect to root 'index' */

			case "stop":

				threads.Thread{}.Terminate(fh.sessionData, "") /* Terminate ThreadID */

			case "update":

				functions.SaveConfigData(fh.viperData, r, fh.sessionData) /* Save the configuration data */
				http.Redirect(w, r, fmt.Sprintf(r.URL.Path), 301)         /* Redirect to root 'index' */

			case "buy":

				fh.sessionData.ForceBuy = true                    /* Force buy */
				http.Redirect(w, r, fmt.Sprintf(r.URL.Path), 301) /* Redirect to root 'index' */

			case "sell":

				if r.PostFormValue("orderID") == "" { /* Check if the orderID is empty */

					fh.sessionData.ForceSellOrderID = 0 /* Force sell most recent order */
					fh.sessionData.ForceSell = true     /* Force sell */

				} else {

					fh.sessionData.ForceSellOrderID = functions.StrToInt(r.PostFormValue("orderID")) /* Force sell a specific orderID */
					fh.sessionData.ForceSell = true                                                  /* Force sell */

				}

				http.Redirect(w, r, fmt.Sprintf(r.URL.Path), 301) /* Redirect to root 'index' */

			case "configTemplate":

				fh.sessionData.ConfigTemplate = functions.StrToInt(r.PostFormValue("configTemplateList")) /* Retrieve Configuration Template Key selection */
				configData := functions.LoadConfigTemplate(fh.viperData, fh.sessionData)                  /* Load the configuration data */
				functions.ExecuteTemplate(w, configData, fh.sessionData)                                  /* This is the template execution for 'index' */

			}
		}

	}

}

func execution(
	viperData *types.ViperData,
	configData *types.Config,
	sessionData *types.Session,
	marketData *types.Market) {

	var err error /* Error handling */

	/* Connect to Exchange */
	if err = exchange.GetClient(configData, sessionData); err != nil { /* GetClient returns an error if the connection to the exchange is not successful */

		threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error()) /* Terminate ThreadID */

	}

	/* Routine to resume operations */
	var threadIDSessionDB string

	if sessionData.ThreadID, threadIDSessionDB, err = mysql.GetThreadTransactionDistinct(sessionData); err != nil { /* GetThreadTransactionDistinct returns an error if the connection to the database is not successful */

		threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error()) /* Terminate ThreadID */

	}

	if sessionData.ThreadID != "" && !configData.NewSession { /* If ThreadID is not empty and NewSession is false */

		threads.Thread{}.Lock(sessionData) /* Lock thread file */

		configData = functions.GetConfigData(viperData, sessionData) /* Get Config Data */

		if sessionData.Symbol, err = mysql.GetOrderSymbol(sessionData); err != nil { /* GetOrderSymbol returns an error if the connection to the database is not successful */

			threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error()) /* Terminate ThreadID */

		}

		/* Select the symbol coin to be used from sessionData.Symbol */
		if sessionData.SymbolFiat, err = algorithms.ParseSymbolFiat(sessionData); err != nil { /* ParseSymbolFiat returns an error if the symbol is not valid */

			threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error()) /* Terminate ThreadID */

		}

		logger.LogEntry{ /* Log Entry */
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  "Resuming on port " + sessionData.Port,
			LogLevel: "InfoLevel",
		}.Do()

	} else { /* If ThreadID is empty or NewSession is true */

		sessionData.ThreadID = functions.GetThreadID() /* Get ThreadID */

		if !(threads.Thread{}.Lock(sessionData)) { /* Lock thread file */

			os.Exit(1)

		}

		/* Select the symbol coin to be used from Config option */
		sessionData.Symbol = configData.Symbol
		sessionData.SymbolFiat = configData.SymbolFiat

		logger.LogEntry{ /* Log Entry */
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  "Initializing on port " + sessionData.Port,
			LogLevel: "InfoLevel",
		}.Do()

	}

	asyncFunctions(viperData, configData, sessionData) /* Starts async functions that are executed at specific intervals */

	/* Retrieve available fiat funds and update database
	This is only used for retrieving balances for the first time, and is then followed by
	the Websocket routine to retrieve realtime user data  */
	if sessionData.SymbolFiatFunds, err = exchange.GetSymbolFiatFunds( /* GetSymbolFiatFunds returns an error if the connection to the exchange is not successful */
		configData,
		sessionData); err == nil { /* If the connection to the exchange is successful */
		_ = mysql.UpdateSession( /* Update database with available fiat funds */
			configData,
			sessionData)
	}

	/* Retrieve available symbol funds
	This is only used for retrieving balances for the first time, ans is then followed by
	the Websocket routine to retrieve realtime user data  */
	sessionData.SymbolFunds, err = exchange.GetSymbolFunds(configData, sessionData)

	/* Retrieve exchange lot size for ticker and store in sessionData */
	exchange.GetLotSize(configData, sessionData)

	sum := 0
	for {

		/* Check start/stop times of operation */
		if configData.TimeEnforce { /* If TimeEnforce is true */

			for !functions.IsInTimeRange(configData.TimeStart, configData.TimeStop) { /* If current time is not in the time range */

				logger.LogEntry{ /* Log Entry */
					Config:   configData,
					Market:   marketData,
					Session:  sessionData,
					Order:    &types.Order{},
					Message:  "Sleeping",
					LogLevel: "InfoLevel",
				}.Do()

				time.Sleep(300000 * time.Millisecond) /* Sleep for 5 minutes */

			}

		}

		/* Update ThreadCount */
		sessionData.ThreadCount, err = mysql.GetThreadTransactionCount(sessionData)

		/* Update Number of Sale Transactions per hour */
		sessionData.SellTransactionCount, err = mysql.GetOrderTransactionCount(sessionData, "SELL")

		/* This routine is executed when no transaction cycle has initiated (ThreadCount = 0) */
		if sessionData.ThreadCount == 0 { /* If ThreadCount is 0 */

			sessionData.ThreadIDSession = functions.GetThreadID() /* Get ThreadID */

			/* Save new session to Session table. */
			if err := mysql.SaveSession(
				configData,
				sessionData); err != nil {

				/* Update existing session on Session table */
				if err := mysql.UpdateSession(
					configData,
					sessionData); err != nil {

					threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error()) /* Terminate ThreadID */

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
						threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

					}

				}

			}

		}

		/* Conditional used in case this is the first run in the cycle go get past market data */
		if marketData.PriceChangeStatsHighPrice == 0 { /* If PriceChangeStatsHighPrice is 0 */

			markets.Data{}.LoadKlinePast(configData, marketData, sessionData) /* Load Kline Past */

		}

		wg := &sync.WaitGroup{} /* WaitGroup to stop inside Channels */
		wg.Add(3)               /* WaitGroup to stop inside Channels */

		go telegram.CheckUpdates( /* Check for Telegram updates */
			configData,
			sessionData,
			wg)

		go algorithms.WsKline( /* Websocket routine to retrieve realtime candle data */
			configData,
			marketData,
			sessionData,
			wg)

		go algorithms.WsUserDataServe( /* Websocket routine to retrieve realtime user data */
			configData,
			sessionData,
			wg)

		go algorithms.WsBookTicker( /* Websocket routine to retrieve realtime ticker prices */
			viperData,
			configData,
			marketData,
			sessionData,
			wg)

		wg.Wait() /* Wait for the goroutines to finish */

		logger.LogEntry{ /* Log Entry */
			Config:   configData,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  "All websocket channels stopped",
			LogLevel: "DebugLevel",
		}.Do()

		sessionData.StopWs = false /* Reset goroutine channels */

		/* Reload configuration in case of WsBookTicker broken connection */
		configData = functions.GetConfigData(viperData, sessionData) /* Get Config Data */

		time.Sleep(3000 * time.Millisecond) /* Sleep for 3 seconds */

		logger.LogEntry{ /* Log Entry */
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

// asyncFunctions starts async functions that are executed at specific intervals
func asyncFunctions(
	viperData *types.ViperData,
	configData *types.Config,
	sessionData *types.Session) {

	/* Synchronize time with Binance every 5 minutes */
	_ = exchange.NewSetServerTimeService(configData, sessionData)
	scheduler.RunTaskAtInterval(
		func() { _ = exchange.NewSetServerTimeService(configData, sessionData) },
		time.Second*300,
		time.Second*0)

	/* Retrieve config data every 10 seconds. */
	scheduler.RunTaskAtInterval(
		func() { configData = functions.GetConfigData(viperData, sessionData) },
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
	nodes.Node{}.GetRole(configData, sessionData)
	scheduler.RunTaskAtInterval(
		func() {
			nodes.Node{}.GetRole(configData, sessionData)
		},
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
			sessionData.SellTransactionCount, _ = mysql.GetOrderTransactionCount(sessionData, "SELL")
		},
		time.Second*180,
		time.Second*0)

	/* Update exchange latency every 5 seconds. */
	scheduler.RunTaskAtInterval(
		func() {
			sessionData.Latency, _ = functions.GetExchangeLatency(sessionData)
		},
		time.Second*5,
		time.Second*0)

	/* Check system status every 10 seconds. */
	scheduler.RunTaskAtInterval(
		func() {
			nodes.Node{}.CheckStatus(configData, sessionData)
		},
		time.Second*10,
		time.Second*0)

	/* Send Telegram message with system error (only Master Node) every 60 seconds. */
	scheduler.RunTaskAtInterval(
		func() {
			if sessionData.MasterNode && sessionData.TgBotAPIChatID != 0 {
				if threadID, err := mysql.GetSessionStatus(sessionData); err == nil {
					if threadID != "" {
						telegram.Message{
							Text: "\f" + "System Fault @ " + threadID,
						}.Send(sessionData)
					}
				}
			}
		}, time.Second*60,
		time.Second*0)

	/* Load mySQL dynamic components for javascript autoloader every 10 seconds. */
	scheduler.RunTaskAtInterval(
		func() {
			loader.LoadSessionDataAdditionalComponentsAsync(sessionData)
		},
		time.Second*10,
		time.Second*0)

}
