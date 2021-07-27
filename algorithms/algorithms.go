package algorithms

import (
	"cryptopump/exchange"
	"cryptopump/functions"
	"cryptopump/markets"
	"cryptopump/mysql"
	"cryptopump/plotter"
	"cryptopump/threads"
	"cryptopump/types"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2"

	log "github.com/sirupsen/logrus"
)

/* Modify profit based on sell transaction count  */
func calculateProfit(
	configData *types.Config,
	sessionData *types.Session) (profit float64) {

	profit = functions.StrToFloat64(configData.Profit_min.(string))

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
			orderStatus.CummulativeQuoteQuantity,
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

	return marketData.Price >= (marketData.PriceChangeStatsHighPrice * (1 - functions.StrToFloat64(configData.Buy_24hs_highprice_entry.(string))))

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
	if functions.StrToFloat64(configData.Buy_quantity_fiat_up.(string)) == 0 {

		return false, 0

	}

	/* Validate RSI7 lower than buy_rsi7_entry */
	if marketData.Rsi7 > functions.StrToFloat64(configData.Buy_rsi7_entry.(string)) {

		return false, 0

	}

	/* If Market Direction is less than configData.Buy_direction_up do not buy. Defined in WsKline. */
	if marketData.Direction < functions.StrToInt(configData.Buy_direction_up.(string)) {

		return false, 0

	}

	if lastOrderTransactionPrice, err = mysql.GetLastOrderTransactionPrice(
		sessionData,
		"SELL"); err != nil {

		return false, 0

	}

	/* Test if event price is lower than last Sell price plus threshold up */
	buyRepeatThresholdUp := functions.StrToFloat64(configData.Buy_repeat_threshold_up.(string))
	if marketData.Price < lastOrderTransactionPrice*(1+buyRepeatThresholdUp) {

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
		order.CummulativeQuoteQuantity,
		order.TransactTime,
		err = mysql.GetThreadLastTransaction(sessionData); err != nil {

		return false, 0

	}

	/* See comment above */
	if marketData.Price > order.Price &&
		marketData.Price < (order.Price*(1+(functions.StrToFloat64(configData.Profit_min.(string))/2))) {

		return false, 0

	} else if marketData.Price < order.Price &&
		marketData.Price > (order.Price*(1-(functions.StrToFloat64(configData.Profit_min.(string))/2))) {

		return false, 0

	}

	/* 		This function retrieve the number of thread transactions with price bigger than current price times buy_repeat_threshold_up.
	   		It servers the purpose of ensuring the algorithm does not buy above the biggest buy. If more more than 1 transaction will not execute buy. */
	if threadTransactiontUpmarketPriceCount, err = mysql.GetThreadTransactiontUpmarketPriceCount(
		sessionData,
		(marketData.Price * (1 + buyRepeatThresholdUp))); err != nil {

		return false, 0

	}

	/* See comment above */
	if functions.IntToFloat64(threadTransactiontUpmarketPriceCount) > 1 {

		return false, 0

	}

	functions.Logger(
		configData,
		marketData,
		sessionData,
		log.InfoLevel,
		0,
		0,
		0,
		0,
		"UP")

	switch {
	case sessionData.ThreadCount == 1:

		/* Stop  large transactions at the top os the order book. */
		return true, functions.StrToFloat64(configData.Buy_quantity_fiat_init.(string))

	case sessionData.ThreadCount > functions.StrToInt(configData.Buy_repeat_threshold_down_second_start_count.(string)):

		/* Stop large transactions if count is bigger than specified in config. */
		return true, functions.StrToFloat64(configData.Buy_quantity_fiat_init.(string))

	default:

		return true, functions.StrToFloat64(configData.Buy_quantity_fiat_up.(string))

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
	if configData.Buy_quantity_fiat_down == 0 {

		return false, 0

	}

	/* Validate RSI14 not negative */
	if marketData.Rsi14 <= 0 {

		return false, 0

	}

	/* Validate market direction is uptrend */
	if marketData.Direction < functions.StrToInt(configData.Buy_direction_down.(string)) {

		return false, 0

	}

	/* Ensure funds are not deployed less than buy_repeat_threshold_down from each other */
	buyRepeatThresholdDown := functions.StrToFloat64(configData.Buy_repeat_threshold_down.(string))
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

		buyRepeatThresholdDown = functions.StrToFloat64(configData.Buy_repeat_threshold_down_second.(string))

	}

	/* Test with new buy_repeat_threshold_down */
	if marketData.Price > (lastOrderTransactionPrice * (1 - buyRepeatThresholdDown)) {

		return false, 0

	}

	functions.Logger(
		configData,
		marketData,
		sessionData,
		log.InfoLevel,
		0,
		0,
		0,
		0,
		"DOWN")

	return true, functions.StrToFloat64(configData.Buy_quantity_fiat_down.(string))

}

func isBuyInitial(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) (bool, float64) {

	/* Validate RSI7 lower than buy_rsi7_entry */
	/* Validate RSI3 not negative */
	if marketData.Rsi7 < functions.StrToFloat64(configData.Buy_rsi7_entry.(string)) && marketData.Rsi3 > 0 {

		/* Do not log if DryRun mode set to true */
		if configData.DryRun != "true" {

			functions.Logger(
				configData,
				marketData,
				sessionData,
				log.InfoLevel,
				0,
				0,
				0,
				0,
				"INIT")

		}

		return true, functions.StrToFloat64(configData.Buy_quantity_fiat_init.(string))

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
	sessionData.ListenKey, _ = exchange.GetUserStreamServiceListenKey(configData, sessionData)

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

		} else if executionReport.EventType == "executionReport" {

			return

		}

		/* Unmarshal and process outboundAccountPosition */
		if err := json.Unmarshal(message, &outboundAccountPosition); err != nil {

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

		} else if outboundAccountPosition.EventType == "outboundAccountPosition" {

			for key := range outboundAccountPosition.Balances {

				if outboundAccountPosition.Balances[key].Asset == sessionData.Symbol_fiat {

					sessionData.Symbol_fiat_funds = functions.StrToFloat64(outboundAccountPosition.Balances[key].Free)

					_ = mysql.UpdateSession(
						configData,
						sessionData)

				}

			}

			return

		}

	}

	errHandler := func(err error) {

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

		stopChannels(stopC, wg, configData, sessionData)

		/* Retrieve NEW WsUserDataServe listen key for user stream service when there's an error */
		sessionData.ListenKey, _ = exchange.GetUserStreamServiceListenKey(configData, sessionData)

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
				exchange.BinanceMapWsKline(event.Kline))

		}

	}

	errHandler := func(err error) {

		functions.Logger(
			configData,
			marketData,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

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
			configData.Exit.(string) == "true" {

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
				sessionData.ThreadCount, _ = mysql.GetThreadTransactionCount(sessionData)

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

		functions.Logger(
			configData,
			marketData,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

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

	/* Validate available funds to buy */
	if !functions.IsFundsAvailable(
		configData,
		sessionData) {

		return false, 0

	}

	/* Trigger Force Buy */
	if sessionData.ForceBuy {

		sessionData.ForceBuy = false

		return true, functions.StrToFloat64(configData.Buy_quantity_fiat_init.(string))

	}

	/* If configData.Exit is True stop BUY. */
	if configData.Exit.(string) == "true" {

		return false, 0

	}

	/* Validate marketData not older than 100 seconds */
	if time.Since(marketData.TimeStamp).Seconds() > 100 {

		return false, 0

	}

	/* 	If last buy is less than configData.Buy_wait seconds return false
	   	This function protects against sequential buys when there's too much volatility */
	if time.Duration(time.Since(sessionData.LastBuyTransactTime).Seconds()) < time.Duration(functions.StrToFloat64(configData.Buy_wait.(string))) {

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
			order.CummulativeQuoteQuantity,
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
	if time.Duration(time.Since(sessionData.LastSellCanceledTime).Seconds()) < time.Duration(functions.StrToInt(configData.SellWaitAfterCancel.(string))) {

		return false, order

	}

	/* Sell-to-Cover - Sell if Fiat funds are lower than buy qty and ticker price is below last thread transaction.
	This will sell at loss, but make funds available for new buy transactions */
	if configData.Exit.(string) != "true" && /* Doesn't force sell if system is in Exit mode */
		configData.SellToCover.(string) == "true" { /* Doesn't force sell if SellToCover is False */

		if (sessionData.Symbol_fiat_funds - functions.StrToFloat64(configData.Symbol_fiat_stash.(string))) < functions.StrToFloat64(configData.Buy_quantity_fiat_down.(string)) {

			/* Retrieve the last 'active' BUY transaction for a Thread */
			order.OrderID,
				order.Price,
				order.ExecutedQuantity,
				order.CummulativeQuoteQuantity,
				order.TransactTime,
				_ = mysql.GetThreadLastTransaction(sessionData)

			if marketData.Price < (order.Price * (1 - functions.StrToFloat64(configData.Buy_repeat_threshold_down.(string)))) {

				return true, order

			}

		}
	}

	/* Retrieve lowest price order from Thread database */
	if order.OrderID,
		order.Price,
		order.ExecutedQuantity,
		order.CummulativeQuoteQuantity,
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

	/* Verify that an order is in a sellable time range
	This function help to avoid issue when a sale happen in the same second as the Buy transaction.
	Duration must be provided in seconds */
	if !isOrderInTimeRangeToSell(order, 60) {

		return false, order

	}

	/* Current price is higher than BUY price + profits */
	/* Modify profit based on sell transaction count  */
	if (marketData.Price*(1+functions.StrToFloat64(configData.Exchange_comission.(string)))) >=
		(order.Price*(1+calculateProfit(configData, sessionData))) &&
		order.OrderID != 0 {

		/* Hold sale if RSI3 above defined threshold.
		The objective of this setting is to extend the holding as long as possible while ticker price is climbing */
		if marketData.Rsi3 > functions.StrToFloat64(configData.SellHoldOnRSI3.(string)) {

			return false, order

		}

		return true, order

	}

	return false, order

}
