package exchange

import (
	"math"
	"strings"
	"time"

	"cryptopump/functions"
	"cryptopump/mysql"
	"cryptopump/threads"
	"cryptopump/types"

	log "github.com/sirupsen/logrus"
)

/* Define the exchange to be used */
func GetClient(
	configData *types.Config,
	sessionData *types.Session) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		sessionData.Clients.Binance = binanceGetClient(configData)

	}

}

/* Retrieve Order Status */
func GetOrder(
	configData *types.Config,
	sessionData *types.Session,
	orderID int64) (order *types.Order, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceGetOrder(sessionData, orderID)

	}

	return

}

/* Create order to BUY */
func BuyOrder(
	configData *types.Config,
	sessionData *types.Session,
	quantity string) (order *types.Order, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceBuyOrder(sessionData, quantity)

	}

	return

}

/* Create order to SELL */
func SellOrder(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session,
	quantity string) (order *types.Order, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceSellOrder(marketData, sessionData, quantity)

	}

	return

}

/* CANCEL an order */
func CancelOrder(
	configData *types.Config,
	sessionData *types.Session,
	orderID int64) (order *types.Order, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceCancelOrder(sessionData, orderID)

	}

	return

}

/* Retrieve exchange information */
func GetInfo(
	configData *types.Config,
	sessionData *types.Session) (info *types.ExchangeInfo, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceGetInfo(sessionData)

	}

	return

}

/* Retrieve Lot Size specs */
func GetLotSize(
	configData *types.Config,
	sessionData *types.Session) {

	if info, err := GetInfo(configData, sessionData); err != nil {

		return

	} else {

		sessionData.MaxQuantity = functions.StrToFloat64(info.MaxQuantity)
		sessionData.MinQuantity = functions.StrToFloat64(info.MinQuantity)
		sessionData.StepSize = functions.StrToFloat64(info.StepSize)

	}

}

/* Retrieve funds available */
func GetSymbolFunds(
	configData *types.Config,
	sessionData *types.Session) (balance float64, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceGetSymbolFunds(sessionData)

	}

	return

}

/* Retrieve KLines via REST API */
func GetKlines(
	configData *types.Config,
	sessionData *types.Session) (klines []*types.Kline, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		tmp, err := binanceGetKlines(sessionData)

		if err == nil {
			return binanceMapKline(tmp), err
		}

		return nil, err

	}

	return

}

/* Retrieve 24hs Rolling Price Statistics */
func GetPriceChangeStats(
	configData *types.Config,
	sessionData *types.Session,
	marketData *types.Market) (priceChangeStats []*types.PriceChangeStats, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceGetPriceChangeStats(sessionData)

	}

	return

}

/* Calculate the correct quantity to SELL according to the exchange lotSizeStep */
func getSellQuantity(
	order types.Order,
	sessionData *types.Session) (quantity float64) {

	return math.Round(order.ExecutedQuantity/sessionData.StepSize) * sessionData.StepSize

}

/* Calculate the correct quantity to BUY according to the exchange lotSizeMin and lotSizeStep */
func getBuyQuantity(
	marketData *types.Market,
	sessionData *types.Session,
	fiatQuantity float64) (quantity float64) {

	return functions.ConvertFiatToCoin(fiatQuantity, marketData.Price, sessionData.MinQuantity, sessionData.StepSize)

}

/* Retrieve listen key for user stream service */
func GetUserStreamServiceListenKey(
	configData *types.Config,
	sessionData *types.Session) (listenKey string, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceGetUserStreamServiceListenKey(sessionData)

	}

	return

}

/* Keep user stream service alive */
func KeepAliveUserStreamServiceListenKey(
	configData *types.Config,
	sessionData *types.Session) (err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceKeepAliveUserStreamServiceListenKey(sessionData)

	}

	return

}

/* Synchronize time */
func NewSetServerTimeService(
	configData *types.Config,
	sessionData *types.Session) (err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceNewSetServerTimeService(sessionData)

	}

	return

}

/* WsBookTickerServe serve websocket that pushes updates to the best bid or ask price or quantity in real-time for a specified symbol. */
func WsBookTickerServe(
	configData *types.Config,
	sessionData *types.Session,
	wsHandler *types.WsHandler,
	errHandler func(err error)) (doneC chan struct{}, stopC chan struct{}, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceWsBookTickerServe(sessionData, wsHandler, errHandler)

	}

	return

}

/* WsKlineServe serve websocket kline handler */
func WsKlineServe(
	configData *types.Config,
	sessionData *types.Session,
	wsHandler *types.WsHandler,
	errHandler func(err error)) (doneC chan struct{}, stopC chan struct{}, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceWsKlineServe(sessionData, wsHandler, errHandler)

	}

	return

}

/* WsUserDataServe serve user data handler with listen key */
func WsUserDataServe(
	configData *types.Config,
	sessionData *types.Session,
	wsHandler *types.WsHandler,
	errHandler func(err error)) (doneC chan struct{}, stopC chan struct{}, err error) {

	switch strings.ToLower(configData.ExchangeName.(string)) {
	case "binance":

		return binanceWsUserDataServe(sessionData, wsHandler, errHandler)

	}

	return

}

/* Buy Ticker */
func BuyTicker(
	quantity float64,
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) {

	var orderStatus *types.Order
	var orderPrice float64
	var orderExecutedQuantity float64
	var isCanceled bool

	/* Enter and defer exiting busy mode */
	sessionData.Busy = true
	defer func() {
		sessionData.Busy = false
	}()

	/* Exit if DryRun mode set to true */
	if configData.DryRun == "true" {

		return

	}

	orderResponse, err := BuyOrder(
		configData,
		sessionData,
		functions.Float64ToStr(getBuyQuantity(marketData, sessionData, quantity), 4)) /* Get the correct quantity according to lotSizeMin and lotSizeStep */

	/* Test orderResponse for  errors */
	if (orderResponse == nil && err != nil) ||
		(orderResponse == nil && err == nil) {

		return

	}

	/* Check if result is nil and set as zero */
	if orderPrice = orderResponse.CummulativeQuoteQuantity / orderResponse.ExecutedQuantity; math.IsNaN(orderPrice) {
		orderPrice = 0
	}

	orderExecutedQuantity = orderResponse.ExecutedQuantity

	/* Save order to database */
	if err := mysql.SaveOrder(
		sessionData,
		orderResponse.ClientOrderID,
		orderResponse.CummulativeQuoteQuantity,
		orderResponse.ExecutedQuantity,
		int64(orderResponse.OrderID),
		orderPrice,
		string(orderResponse.Side),
		string(orderResponse.Status),
		orderResponse.Symbol,
		orderResponse.TransactTime); err != nil {

		/* Cleanly exit ThreadID */
		threads.ExitThreadID(sessionData)

	}

	/* This session variable stores the time of the last buy */
	sessionData.LastBuyTransactTime = time.Now()

S:
	switch orderResponse.Status {
	case "FILLED", "PARTIALLY_FILLED":
	case "CANCELED":

		isCanceled = true

	case "NEW":

		for orderStatus, err = GetOrder(
			configData,
			sessionData,
			int64(orderResponse.OrderID)); orderStatus == nil || orderStatus.Status == "NEW"; {

			if err != nil {

				break S

			}

			time.Sleep(3000 * time.Millisecond)

		}

		switch orderStatus.Status {
		case "FILLED", "PARTIALLY_FILLED":

			orderPrice = orderStatus.CummulativeQuoteQuantity / orderStatus.ExecutedQuantity

			orderExecutedQuantity = orderStatus.ExecutedQuantity

			/* Update order status and price & Save Thread Transaction */
			if err := mysql.UpdateOrder(
				sessionData,
				int64(orderResponse.OrderID),
				orderResponse.CummulativeQuoteQuantity,
				orderResponse.ExecutedQuantity,
				orderPrice,
				string(orderStatus.Status)); err != nil {

				/* Cleanly exit ThreadID */
				threads.ExitThreadID(sessionData)

			}

		case "CANCELED":

			isCanceled = true

			break S

		}

	}

	if !isCanceled {

		/* Save Thread Transaction */
		if err := mysql.SaveThreadTransaction(
			sessionData,
			int64(orderResponse.OrderID),
			orderResponse.CummulativeQuoteQuantity,
			orderPrice,
			orderExecutedQuantity); err != nil {

			/* Cleanly exit ThreadID */
			threads.ExitThreadID(sessionData)

		}

		functions.Logger(
			configData,
			nil,
			sessionData,
			log.InfoLevel,
			0,
			int64(orderResponse.OrderID),
			orderPrice,
			0,
			"BUY")

	} else if isCanceled {

		functions.Logger(
			configData,
			nil,
			sessionData,
			log.InfoLevel,
			0,
			int64(orderResponse.OrderID),
			orderPrice,
			0,
			"CANCELED")

	}

}

/* Sell Ticker */
func SellTicker(
	order types.Order,
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) {

	var orderResponse *types.Order
	var orderStatus *types.Order
	var cummulativeQuoteQuantity float64
	var cancelOrderResponse *types.Order
	var isCanceled bool
	var err error
	var i int

	/* Enter and defer exiting busy mode */
	sessionData.Busy = true
	defer func() {
		sessionData.Busy = false
	}()

	/* Exit if DryRun mode set to true */
	if configData.DryRun == "true" {

		return

	}

	orderResponse, err = SellOrder(
		configData,
		marketData,
		sessionData,
		functions.Float64ToStr(getSellQuantity(order, sessionData), 6) /* Get correct quantity to sell according to the lotSizeStep */)

	/* Test orderResponse for  errors */
	if (orderResponse == nil && err != nil) ||
		(orderResponse == nil && err == nil) {

		return

	}

	/* Save order to database */
	if err := mysql.SaveOrder(
		sessionData,
		orderResponse.ClientOrderID,
		orderResponse.CummulativeQuoteQuantity,
		orderResponse.ExecutedQuantity,
		int64(orderResponse.OrderID),
		marketData.Price,
		string(orderResponse.Side),
		string(orderResponse.Status),
		orderResponse.Symbol,
		orderResponse.TransactTime); err != nil {

		/* Cleanly exit ThreadID */
		threads.ExitThreadID(sessionData)

	}

S:
	switch orderResponse.Status {
	case "FILLED":

		cummulativeQuoteQuantity = orderResponse.CummulativeQuoteQuantity

	case "CANCELED":

		isCanceled = true

	case "PARTIALLY_FILLED", "NEW":

		time.Sleep(2000 * time.Millisecond)

	F:
		for orderStatus, err = GetOrder(
			configData,
			sessionData,
			int64(orderResponse.OrderID)); orderStatus == nil ||
			orderStatus.Status == "NEW" ||
			orderStatus.Status == "PARTIALLY_FILLED"; {

			if err != nil {

				/* Cleanly exit ThreadID */
				threads.ExitThreadID(sessionData)

			}

			switch orderStatus.Status {
			case "FILLED":

				break F

			case "CANCELED":

				isCanceled = true

				break F

			}

			i++ /* increment iterations before order cancel */

			/* Initiate order cancel after 10 iterations */
			if i == 9 {

				if cancelOrderResponse, err = CancelOrder(
					configData,
					sessionData,
					int64(orderResponse.OrderID)); err != nil {

					switch {
					case strings.Contains(err.Error(), "-2010"), strings.Contains(err.Error(), "-2011"):
						/* -2011 Order filled in full before cancelling */
						/* -2010 Account has insufficient balance for requested action */

						if orderStatus, err = GetOrder(
							configData,
							sessionData,
							int64(orderResponse.OrderID)); err != nil {

							/* Cleanly exit ThreadID */
							threads.ExitThreadID(sessionData)

						}

						break F

					default:

						functions.Logger(
							configData,
							nil,
							sessionData,
							log.DebugLevel,
							0,
							int64(orderResponse.OrderID),
							0,
							0,
							err.Error())

						break S
					}

				}

				switch cancelOrderResponse.Status {
				case "CANCELED":

					isCanceled = true

					/* This session variable stores the time of the cancelled sell */
					sessionData.LastSellCanceledTime = time.Now()

					if orderStatus, err = GetOrder(
						configData,
						sessionData,
						int64(orderResponse.OrderID)); err != nil {

						/* Cleanly exit ThreadID */
						threads.ExitThreadID(sessionData)

					}

					break F

				default:

					functions.Logger(
						configData,
						nil,
						sessionData,
						log.InfoLevel,
						order.OrderID,
						int64(orderResponse.OrderID),
						marketData.Price,
						0,
						"FAILED TO CANCEL ORDER")

					break F

				}

			}

			/* Wait time between iterations (i++). There are ten iterations and the total waiting time define the amount od time before an order is canceled. configData.SellWaitBeforeCancel is divided by then converted into seconds. */
			time.Sleep(
				time.Duration(
					functions.StrToInt(configData.SellWaitBeforeCancel.(string))/10) * time.Second)

		}

		cummulativeQuoteQuantity = orderStatus.CummulativeQuoteQuantity

		/* Update order status and price */
		if err := mysql.UpdateOrder(
			sessionData,
			int64(orderResponse.OrderID),
			orderStatus.CummulativeQuoteQuantity,
			orderStatus.ExecutedQuantity,
			marketData.Price,
			string(orderStatus.Status)); err != nil {

			/* Cleanly exit ThreadID */
			threads.ExitThreadID(sessionData)

		}

	}

	if !isCanceled {

		/* Remove Thread transaction from database */
		if err := mysql.DeleteThreadTransactionByOrderID(
			sessionData,
			order.OrderID); err != nil {

			/* Cleanly exit ThreadID */
			threads.ExitThreadID(sessionData)

		}

		functions.Logger(
			configData,
			nil,
			sessionData,
			log.InfoLevel,
			order.OrderID,
			int64(orderResponse.OrderID),
			marketData.Price,
			functions.GetProfitResult(order.CummulativeQuoteQuantity, cummulativeQuoteQuantity),
			"SELL")

	} else if isCanceled {

		functions.Logger(
			configData,
			nil,
			sessionData,
			log.InfoLevel,
			order.OrderID,
			int64(orderResponse.OrderID),
			marketData.Price,
			0,
			"CANCELED")

	}

}
