package loader

/* This package contains the functions responsible for loading frequent data
that is made available to the javascript autoloader for html output. This data
is commonly loaded via the webserver using GET/sessiondata */

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/mysql"
	"github.com/aleibovici/cryptopump/types"
)

// LoadSessionDataAdditionalComponents Load dynamic components for javascript autoloader for html output
func LoadSessionDataAdditionalComponents(
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
		ProfitThreadIDPct    float64 /* ThreadID profit percentage */
		Profit               float64 /* Total profit */
		ProfitPct            float64 /* Total profit percentage */
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

	sessiondata.Session.Profit = math.Round(sessionData.Global.Profit*100) / 100                       /* Sessions.Global loaded from mySQL via loadSessionDataAdditionalComponentsAsync */
	sessiondata.Session.ProfitPct = math.Round(sessionData.Global.ProfitPct*100) / 100                 /* Sessions.Global loaded from mySQL via loadSessionDataAdditionalComponentsAsync */
	sessiondata.Session.ProfitThreadID = math.Round(sessionData.Global.ProfitThreadID*100) / 100       /* Sessions.Global loaded from mySQL via loadSessionDataAdditionalComponentsAsync */
	sessiondata.Session.ProfitThreadIDPct = math.Round(sessionData.Global.ProfitThreadIDPct*100) / 100 /* Sessions.Global loaded from mySQL via loadSessionDataAdditionalComponentsAsync */
	sessiondata.Session.ThreadCount = sessionData.Global.ThreadCount                                   /* Sessions.Global loaded from mySQL via loadSessionDataAdditionalComponentsAsync */
	sessiondata.Session.ThreadAmount = math.Round(sessionData.Global.ThreadAmount*100) / 100           /* Sessions.Global loaded from mySQL via loadSessionDataAdditionalComponentsAsync */

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

// LoadSessionDataAdditionalComponentsAsync Load mySQL dynamic components for javascript autoloader for html output.
// This is a separate function because it is reload with scheduler.RunTaskAtInterval via asyncFunctions
func LoadSessionDataAdditionalComponentsAsync(sessionData *types.Session) {

	var err error

	/* Conditional defer logging when there is an error retriving data */
	defer func() {
		if err != nil {
			logger.LogEntry{
				Config:   nil,
				Market:   nil,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "DebugLevel",
			}.Do()
		}
	}()

	/* Load total profit and total profit percentage  */
	if sessionData.Global.Profit, sessionData.Global.ProfitPct, err = mysql.GetProfit(sessionData); err != nil {

		return

	}

	/* Load total thread profit and total thread profit percentage  */
	if sessionData.Global.ProfitThreadID, sessionData.Global.ProfitThreadIDPct, err = mysql.GetProfitByThreadID(sessionData); err != nil {

		return

	}

	/* Load running thread count */
	if sessionData.Global.ThreadCount, err = mysql.GetThreadCount(sessionData); err != nil {

		return

	}

	/* Load total thread dollar amount */
	if sessionData.Global.ThreadAmount, err = mysql.GetThreadAmount(sessionData); err != nil {

		return

	}

}
