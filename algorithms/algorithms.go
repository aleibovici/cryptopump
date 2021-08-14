package algorithms

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aleibovici/cryptopump/exchange"
	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/markets"
	"github.com/aleibovici/cryptopump/mysql"
	"github.com/aleibovici/cryptopump/plotter"
	"github.com/aleibovici/cryptopump/threads"
	"github.com/aleibovici/cryptopump/types"

	"github.com/adshao/go-binance/v2"
)

/* Modify profit based on sell transaction count  */
func calculateProfit(
	configData *types.Config,
	sessionData *types.Session) (profit float64) {

	profit = configData.ProfitMin

	switch {
	case sessionData.SellTransactionCount <= 2:

		profit *= 1

	case sessionData.SellTransactionCount <= 3:

		profit *= 2

	case sessionData.SellTransactionCount > 3:

		profit *= 2.5

	}

	return profit

}

// UpdatePendingOrders Routine to fill rogue and not up-to-date orders in the db and update
func UpdatePendingOrders(
	configData *types.Config,
	sessionData *types.Session) {

	var err error
	var orderID int64
	var orderStatus *types.Order

	if orderID, _, err = mysql.GetOrderTransactionPending(sessionData); err != nil {

		/* Cleanly exit ThreadID */
		threads.ExitThreadID(sessionData)

	}

	if orderID != 0 {

		if orderStatus, err = exchange.GetOrder(
			configData,
			sessionData,
			orderID); err != nil {

			return

		}

		/* Update order status */
		if err := mysql.UpdateOrder(
			sessionData,
			int64(orderStatus.OrderID),
			orderStatus.CumulativeQuoteQuantity,
			orderStatus.ExecutedQuantity,
			orderStatus.Price,
			string(orderStatus.Status)); err != nil {

			/* Cleanly exit ThreadID */
			threads.ExitThreadID(sessionData)

		}

	}

}

/* Check if ticker price lower than 24hs high price */
func is24hsHighPrice(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) bool {

	return marketData.Price >= (marketData.PriceChangeStatsHighPrice * (1 - configData.Buy24hsHighpriceEntry))

}

/* Verify that an order is in a sellable time range
This function help to avoid issue when a sale happen in the same seccond as the Buy transaction. Duration must be provided in seconds */
func isOrderInTimeRangeToSell(
	order types.Order,
	timeRange time.Duration) bool {

	timeNow := time.Now()
	timeTransaction := time.Unix(order.TransactTime/1000, 0)

	return timeNow.Sub(timeTransaction).Seconds() > float64(timeRange)

}

/* Buy Upmarket algorithms */
func isBuyUpmarket(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) (bool, float64) {

	var err error
	var lastOrderTransactionPrice float64
	var lastOrderTransactionSide string
	var threadTransactiontUpmarketPriceCount int
	var order types.Order

	/* If BUY UP amount is 0 do not buy */
	if configData.BuyQuantityFiatUp == 0 {

		return false, 0

	}

	/* Validate RSI7 lower than buy_rsi7_entry */
	if marketData.Rsi7 > configData.BuyRsi7Entry {

		return false, 0

	}

	/* If Market Direction is less than configData.BuyDirectionUp do not buy. Defined in WsKline. */
	if marketData.Direction < configData.BuyDirectionUp {

		return false, 0

	}

	if lastOrderTransactionPrice, err = mysql.GetLastOrderTransactionPrice(
		sessionData,
		"SELL"); err != nil {

		return false, 0

	}

	/* Test if event price is lower than last Sell price plus threshold up */
	if marketData.Price < lastOrderTransactionPrice*(1+configData.BuyRepeatThresholdUp) {

		return false, 0

	}

	/* Retrieve the last transaction side and if it's a BUY exit.
	This avoid double BUY on the UP side */
	if lastOrderTransactionSide, err = mysql.GetLastOrderTransactionSide(sessionData); err != nil {

		return false, 0

	}

	/* Avoid double BUY in UpMarket. Lowest price transaction must be sold first. */
	if lastOrderTransactionSide == "BUY" {

		return false, 0

	}

	/* 	This function retrieve the next transaction from Thread database and verify that
	the ticker price not half profit close to the transaction.This function avoid multiple
	upmarket buy close to each other. */
	if order.OrderID,
		order.Price,
		order.ExecutedQuantity,
		order.CumulativeQuoteQuantity,
		order.TransactTime,
		err = mysql.GetThreadLastTransaction(sessionData); err != nil {

		return false, 0

	}

	/* See comment above */
	if marketData.Price > order.Price &&
		marketData.Price < (order.Price*(1+(configData.ProfitMin/2))) {

		return false, 0

	} else if marketData.Price < order.Price &&
		marketData.Price > (order.Price*(1-(configData.ProfitMin/2))) {

		return false, 0

	}

	/* 		This function retrieve the number of thread transactions with price bigger than current price times buy_repeat_threshold_up.
	   		It servers the purpose of ensuring the algorithm does not buy above the biggest buy. If more more than 1 transaction will not execute buy. */
	if threadTransactiontUpmarketPriceCount, err = mysql.GetThreadTransactiontUpmarketPriceCount(
		sessionData,
		(marketData.Price * (1 + configData.BuyRepeatThresholdUp))); err != nil {

		return false, 0

	}

	/* See comment above */
	if functions.IntToFloat64(threadTransactiontUpmarketPriceCount) > 1 {

		return false, 0

	}

	logger.LogEntry{
		Config:   configData,
		Market:   marketData,
		Session:  sessionData,
		Order:    &types.Order{},
		Message:  "UP",
		LogLevel: "InfoLevel",
	}.Do()

	switch {
	case sessionData.ThreadCount == 1:

		/* Stop  large transactions at the top os the order book. */
		return true, configData.BuyQuantityFiatInit

	case sessionData.ThreadCount > configData.BuyRepeatThresholdDownSecondStartCount:

		/* Stop large transactions if count is bigger than specified in config. */
		return true, configData.BuyQuantityFiatInit

	default:

		return true, configData.BuyQuantityFiatUp

	}

}

/* Buy Downmarket algorithms */
func isBuyDownmarket(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) (bool, float64) {

	var err error
	var lastOrderTransactionPrice float64
	var side1, side2 string

	/* If BUY Down amount is 0 do not buy */
	if configData.BuyQuantityFiatDown == 0 {

		return false, 0

	}

	/* Validate RSI14 not negative */
	if marketData.Rsi14 <= 0 {

		return false, 0

	}

	/* Validate market direction is uptrend */
	if marketData.Direction < configData.BuyDirectionDown {

		return false, 0

	}

	/* Ensure funds are not deployed less than buy_repeat_threshold_down from each other */
	buyRepeatThresholdDown := configData.BuyRepeatThresholdDown
	if lastOrderTransactionPrice, err = mysql.GetLastOrderTransactionPrice(
		sessionData,
		"BUY"); err != nil {

		return false, 0

	}

	/* Test with with buy_repeat_threshold_down to reduce sql queries */
	if marketData.Price > (lastOrderTransactionPrice * (1 - buyRepeatThresholdDown)) {

		return false, 0

	}

	/* Change percentage if last and 2nd orders are BUY */
	if side1, side2, err = mysql.GetOrderTransactionSideLastTwo(sessionData); err != nil {

		return false, 0

	}

	if side1 == "BUY" &&
		side2 == "BUY" {

		buyRepeatThresholdDown = configData.BuyRepeatThresholdDownSecond

	}

	/* Test with new buy_repeat_threshold_down */
	if marketData.Price > (lastOrderTransactionPrice * (1 - buyRepeatThresholdDown)) {

		return false, 0

	}

	logger.LogEntry{
		Config:   configData,
		Market:   marketData,
		Session:  sessionData,
		Order:    &types.Order{},
		Message:  "DOWN",
		LogLevel: "InfoLevel",
	}.Do()

	return true, configData.BuyQuantityFiatDown

}

func isBuyInitial(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) (bool, float64) {

	/* Validate RSI7 lower than buy_rsi7_entry */
	/* Validate RSI3 not negative */
	if marketData.Rsi7 < configData.BuyRsi7Entry && marketData.Rsi3 > 0 {

		/* Do not log if DryRun mode set to true */
		if !configData.DryRun {

			logger.LogEntry{
				Config:   configData,
				Market:   marketData,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  "INIT",
				LogLevel: "InfoLevel",
			}.Do()

		}

		return true, configData.BuyQuantityFiatInit

	}

	return false, 0

}

/* Stop goroutine channels */
func stopChannels(
	channel chan struct{},
	wg *sync.WaitGroup,
	configData *types.Config,
	sessionData *types.Session) {

	sessionData.StopWs = true /* Set goroutine channels to stop */
	channel <- struct{}{}     /* Stop channel that caused initial error */

}

// WsUserDataServe Websocket routine to retrieve realtime user data
func WsUserDataServe(
	configData *types.Config,
	sessionData *types.Session,
	wg *sync.WaitGroup) {

	var doneC chan struct{}
	var stopC chan struct{}
	var err error

	/* Retrieve listen key for user stream service */
	if sessionData.ListenKey, err = exchange.GetUserStreamServiceListenKey(configData, sessionData); err != nil {

		logger.LogEntry{
			Config:   configData,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	wsHandler := &types.WsHandler{}
	wsHandler.BinanceWsUserDataServe = func(message []byte) {

		/* Stop Ws channel */
		if sessionData.StopWs {

			defer wg.Done()
			stopC <- struct{}{}
			return

		}

		var executionReport = &types.ExecutionReport{}
		var outboundAccountPosition = &types.OutboundAccountPosition{}

		/* Unmarshal and process executionReport */
		if err := json.Unmarshal(message, &executionReport); err != nil {

			logger.LogEntry{
				Config:   configData,
				Market:   nil,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "InfoLevel",
			}.Do()

		} else if executionReport.EventType == "executionReport" {

			return

		}

		/* Unmarshal and process outboundAccountPosition */
		if err := json.Unmarshal(message, &outboundAccountPosition); err != nil {

			logger.LogEntry{
				Config:   configData,
				Market:   nil,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "InfoLevel",
			}.Do()

		} else if outboundAccountPosition.EventType == "outboundAccountPosition" {

			for key := range outboundAccountPosition.Balances {

				if outboundAccountPosition.Balances[key].Asset == sessionData.SymbolFiat {

					sessionData.SymbolFiatFunds = functions.StrToFloat64(outboundAccountPosition.Balances[key].Free)

					_ = mysql.UpdateSession(
						configData,
						sessionData)

				}

				/* Update Available crypto funds in exchange */
				if outboundAccountPosition.Balances[key].Asset == sessionData.Symbol[0:3] {

					sessionData.SymbolFunds = functions.StrToFloat64(outboundAccountPosition.Balances[key].Free)

				}

			}

			return

		}

	}

	errHandler := func(err error) {

		logger.LogEntry{
			Config:   configData,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "InfoLevel",
		}.Do()

		switch {
		case strings.Contains(err.Error(), "1001"):
			/* -1001 DISCONNECTED Internal error; unable to process your request. Please try again. */

			exchange.GetClient(configData, sessionData) /* Reconnect exchange client */

		case strings.Contains(err.Error(), "1006"):
			/* -1006 UNEXPECTED_RESP An unexpected response was received from the message bus. Execution status unknown. */
			/* Error Codes for Binance https://github.com/binance/binance-spot-api-docs/blob/master/errors.md */

			return

		case strings.Contains(err.Error(), "read: operation timed out"):
			/* read tcp X.X.X.X:port->X.X.X.X:port: read: operation timed out */

			return

		}

		stopChannels(stopC, wg, configData, sessionData)

		/* Retrieve NEW WsUserDataServe listen key for user stream service when there's an error */
		if sessionData.ListenKey, err = exchange.GetUserStreamServiceListenKey(configData, sessionData); err != nil {

			logger.LogEntry{
				Config:   configData,
				Market:   nil,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "DebugLevel",
			}.Do()

		}

	}

	doneC, stopC, err = exchange.WsUserDataServe(configData, sessionData, wsHandler, errHandler)

	if err != nil {

		fmt.Println(err)

	}

	<-doneC

}

// WsKline The Kline/Candlestick Stream push updates to the current klines/candlestick every second.
func WsKline(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session,
	wg *sync.WaitGroup) {

	var doneC chan struct{}
	var stopC chan struct{}
	var err error

	wsHandler := &types.WsHandler{}
	wsHandler.BinanceWsKline = func(event *binance.WsKlineEvent) {

		/* Stop Ws channel */
		if sessionData.StopWs {

			defer wg.Done()
			stopC <- struct{}{}
			return

		}

		/* Analyse Volume kline direction and create marketData.Direction. 0 = SELL / 1+ BUY */
		activeSellVolume := (functions.StrToFloat64(event.Kline.Volume) - functions.StrToFloat64(event.Kline.ActiveBuyVolume))
		if activeSellVolume > functions.StrToFloat64(event.Kline.ActiveBuyVolume) {

			marketData.Direction = 0

		} else {

			marketData.Direction++

		}

		if event.Kline.IsFinal {

			/* Load Final kline for technical analysis */
			markets.LoadKlineData(
				configData,
				sessionData,
				marketData,
				exchange.BinanceMapWsKline(event.Kline))

			/* Load Final kline for e-chart plotting */
			plotter.LoadKlineData(
				sessionData,
				marketData,
				exchange.BinanceMapWsKline(event.Kline))

		}

	}

	errHandler := func(err error) {

		logger.LogEntry{
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		switch {
		case strings.Contains(err.Error(), "1006"):
			/* -1006 UNEXPECTED_RESP An unexpected response was received from the message bus. Execution status unknown. */
			/* Error Codes for Binance https://github.com/binance/binance-spot-api-docs/blob/master/errors.md */

			return

		case strings.Contains(err.Error(), "unexpected EOF"):
			/* -unexpected EOF An unexpected response was received from the message bus. Execution status unknown. */

			return

		case strings.Contains(err.Error(), "1001"):
			/* -1001 DISCONNECTED Internal error; unable to process your request. Please try again. */

			exchange.GetClient(configData, sessionData) /* Reconnect exchange client */

		}

		stopChannels(stopC, wg, configData, sessionData)

	}

	doneC, stopC, err = exchange.WsKlineServe(configData, sessionData, wsHandler, errHandler)

	if err != nil {

		fmt.Println(err)

	}

	<-doneC

}

// WsBookTicker Pushes any update to the best bid or asks price or quantity in real-time for a specified symbol
func WsBookTicker(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session,
	wg *sync.WaitGroup) {

	var doneC chan struct{}
	var stopC chan struct{}
	var err error

	wsHandler := &types.WsHandler{}
	wsHandler.BinanceWsBookTicker = func(event *binance.WsBookTickerEvent) {

		/* Stop Ws channel */
		if sessionData.StopWs {

			defer wg.Done()
			stopC <- struct{}{}
			return

		}

		/* If there are 0 ThreadID transactions and configData.Exit is True the ThreadID is gracefully
		finalized, and the ThreadID is unlocked. */
		if sessionData.ThreadCount == 0 &&
			configData.Exit {

			/* Delete configuration file for ThreadID */
			functions.DeleteConfigFile(sessionData)

			/* Cleanly exit ThreadID */
			threads.ExitThreadID(sessionData)

		}

		/* Test if event or event.BestAskPrice or marketData are empty or nil before proceeding.
		This test tries to prevent errors where multiple BUYS are executed in a row.
		The source of the problem is unknown but it might be caused by nil data in the event or market data. */
		if event != nil && event.BestAskPrice != "" && marketData != nil {

			marketData.Price = functions.StrToFloat64(event.BestAskPrice)

			if is, buyQuantityFiat := BuyDecisionTree(
				configData,
				marketData,
				sessionData); is {

				exchange.BuyTicker(
					buyQuantityFiat,
					configData,
					marketData,
					sessionData)

				/* Update ThreadCount after BUY */
				sessionData.ThreadCount, err = mysql.GetThreadTransactionCount(sessionData)

			} else if is, order := SellDecisionTree(
				configData,
				marketData,
				sessionData); is {

				exchange.SellTicker(
					order,
					configData,
					marketData,
					sessionData)

				/* Update ThreadCount after SELL */
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

			}

		}

		/* Reload config data every 10 seconds */
		if time.Now().Second()%10 == 0 {

			configData = functions.GetConfigData(sessionData)

		}

	}

	errHandler := func(err error) {

		logger.LogEntry{
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		switch {
		case strings.Contains(err.Error(), "1006"):
			/* -1006 UNEXPECTED_RESP An unexpected response was received from the message bus. Execution status unknown. */
			/* Error Codes for Binance https://github.com/binance/binance-spot-api-docs/blob/master/errors.md */

			return

		case strings.Contains(err.Error(), "1008"):
			/* websocket: close 1008 (policy violation): Pong timeout */
			/* 1008 indicates that an endpoint is terminating the connection
			because it has received a message that violates its policy.  This
			is a generic status code that can be returned when there is no
			other more suitable status code (e.g., 1003 or 1009) or if there
			is a need to hide specific details about the policy. */

			exchange.GetClient(configData, sessionData) /* Reconnect exchange client */

			return

		case strings.Contains(err.Error(), "unexpected EOF"):
			/* -unexpected EOF An unexpected response was received from the message bus. Execution status unknown. */

			return

		case strings.Contains(err.Error(), "connection reset by peer"):

			exchange.GetClient(configData, sessionData) /* Reconnect exchange client */

			return

		case strings.Contains(err.Error(), "read: operation timed out"):
			/* read tcp X.X.X.X:port->X.X.X.X:port: read: operation timed out */

			exchange.GetClient(configData, sessionData) /* Reconnect exchange client */

			return

		}

		stopChannels(stopC, wg, configData, sessionData)

	}

	doneC, stopC, err = exchange.WsBookTickerServe(configData, sessionData, wsHandler, errHandler)

	if err != nil {

		fmt.Println(err)

	}

	<-doneC

}

// BuyDecisionTree BUY decision routine
func BuyDecisionTree(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) (bool, float64) {

	/* Protect against the exchange sending zeroed ticker pricing (seen in few occasions with Binance TestNet)*/
	if marketData.Price == 0 {

		return false, 0

	}

	/* Validate available funds to buy */
	if !functions.IsFundsAvailable(
		configData,
		sessionData) {

		return false, 0

	}

	/* Trigger Force Buy */
	if sessionData.ForceBuy {

		sessionData.ForceBuy = false

		return true, configData.BuyQuantityFiatInit

	}

	/* If configData.Exit is True stop BUY. */
	if configData.Exit {

		return false, 0

	}

	/* Validate marketData not older than 100 seconds */
	if time.Since(marketData.TimeStamp).Seconds() > 100 {

		return false, 0

	}

	/* 	If last buy is less than configData.BuyWait seconds return false
	   	This function protects against sequential buys when there's too much volatility */
	if time.Duration(time.Since(sessionData.LastBuyTransactTime).Seconds()) < time.Duration(configData.BuyWait) {

		return false, 0

	}

	/* Check if ticker price lower than 24hs high price */
	if is24hsHighPrice(
		configData,
		marketData,
		sessionData) {

		return false, 0

	}

	/* Check for subsequent BUY */
	if sessionData.ThreadCount > 0 {

		/* Buy on DOWNMARKET */
		if is, buyQuantityFiat := isBuyDownmarket(
			configData,
			marketData,
			sessionData); is {

			return true, buyQuantityFiat

		}

		/* Buy on UPMARKET */
		if is, buyQuantityFiat := isBuyUpmarket(
			configData,
			marketData,
			sessionData); is {

			return true, buyQuantityFiat

		}

		return false, 0

	}

	/* Check for initial BUY */
	if sessionData.ThreadCount == 0 {

		/* Buy on INITIAL */
		if is, buyQuantityFiat := isBuyInitial(
			configData,
			marketData,
			sessionData); is {

			return true, buyQuantityFiat

		}

		return false, 0

	}

	return false, 0

}

// SellDecisionTree SELL decision routine
func SellDecisionTree(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) (bool, types.Order) {

	var err error
	var order types.Order

	/* Return false if no transactions found */
	if sessionData.ThreadCount == 0 {

		return false, order

	}

	/* Force Sell Most recent open order*/
	if sessionData.ForceSell {

		/* Retrieve the last 'active' BUY transaction for a Thread */
		order.OrderID,
			order.Price,
			order.ExecutedQuantity,
			order.CumulativeQuoteQuantity,
			order.TransactTime,
			_ = mysql.GetThreadLastTransaction(sessionData)

		return true, order

	}

	/* Validate marketData is not older than 100 seconds */
	if time.Since(marketData.TimeStamp).Seconds() > 100 {

		return false, order

	}

	/* 	If last canceled transaction (LastSellCanceledTime) is less than (configData.SellWaitAfterCancel) seconds return false
	   	This function protects against sequential seeling with same pricing */
	if time.Duration(time.Since(sessionData.LastSellCanceledTime).Seconds()) < time.Duration(configData.SellWaitAfterCancel) {

		return false, order

	}

	/* Sell-to-Cover - Sell if Fiat funds are lower than buy qty and ticker price is below last thread transaction.
	This will sell at loss, but make funds available for new buy transactions */
	if !configData.Exit && /* Doesn't force sell if system is in Exit mode */
		configData.SellToCover { /* Doesn't force sell if SellToCover is False */

		if (sessionData.SymbolFiatFunds - configData.SymbolFiatStash) < configData.BuyQuantityFiatDown {

			/* Retrieve the last 'active' BUY transaction for a Thread */
			order.OrderID,
				order.Price,
				order.ExecutedQuantity,
				order.CumulativeQuoteQuantity,
				order.TransactTime,
				_ = mysql.GetThreadLastTransaction(sessionData)

			if marketData.Price < (order.Price * (1 - configData.BuyRepeatThresholdDown)) {

				return true, order

			}

		}
	}

	/* Retrieve lowest price order from Thread database */
	if order.OrderID,
		order.Price,
		order.ExecutedQuantity,
		order.CumulativeQuoteQuantity,
		order.TransactTime,
		err = mysql.GetThreadTransactionByPrice(
		marketData,
		sessionData); err != nil {

		return false, order

	}

	/* If no transactions found return False */
	if order.OrderID == 0 {

		return false, order

	}

	/* Test if symbol funds are available for the Sell order. If not, Buy the amount defined in BuyQuantityFiatInit.
	Sometimes due to decimal changes in transactions or transaction failures there could be divergences and this
	functions help t avoid the problem creating a constant cadence of orders to sell. */
	if sessionData.SymbolFunds <= order.ExecutedQuantity {

		exchange.BuyTicker(
			configData.BuyQuantityFiatInit,
			configData,
			marketData,
			sessionData)

		return false, order

	}

	/* Verify that an order is in a sellable time range
	This function help to avoid issue when a sale happen in the same second as the Buy transaction.
	Duration must be provided in seconds */
	if !isOrderInTimeRangeToSell(order, 60) {

		return false, order

	}

	/* Current price is higher than BUY price + profits */
	/* Modify profit based on sell transaction count  */
	if (marketData.Price*(1+configData.ExchangeComission)) >=
		(order.Price*(1+calculateProfit(configData, sessionData))) &&
		order.OrderID != 0 {

		/* Hold sale if RSI3 above defined threshold.
		The objective of this setting is to extend the holding as long as possible while ticker price is climbing */
		if marketData.Rsi3 > configData.SellHoldOnRSI3 {

			return false, order

		}

		return true, order

	}

	return false, order

}
