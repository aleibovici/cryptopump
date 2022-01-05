package exchange

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/mysql"
	"github.com/aleibovici/cryptopump/threads"
	"github.com/aleibovici/cryptopump/types"
)

// GetClient Define the exchange to be used
func GetClient(
	configData *types.Config,
	sessionData *types.Session) (err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		sessionData.Clients.Binance = binanceGetClient(configData)
		return nil

	}

	return errors.New("Invalid Exchange Name")

}

// GetOrder Retrieve Order Status
func GetOrder(
	configData *types.Config,
	sessionData *types.Session,
	orderID int64) (order *types.Order, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceGetOrder(sessionData, orderID)

	}

	return

}

// BuyOrder Create order to BUY
func BuyOrder(
	configData *types.Config,
	sessionData *types.Session,
	quantity string) (order *types.Order, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceBuyOrder(sessionData, quantity)

	}

	return

}

// SellOrder Create order to SELL
func SellOrder(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session,
	quantity string) (order *types.Order, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceSellOrder(marketData, sessionData, quantity)

	}

	return

}

// CancelOrder CANCEL an order
func CancelOrder(
	configData *types.Config,
	sessionData *types.Session,
	orderID int64) (order *types.Order, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceCancelOrder(sessionData, orderID)

	}

	return

}

// GetInfo Retrieve exchange information
func GetInfo(
	configData *types.Config,
	sessionData *types.Session) (info *types.ExchangeInfo, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceGetInfo(sessionData)

	}

	return

}

// GetLotSize Retrieve Lot Size specs
func GetLotSize(
	configData *types.Config,
	sessionData *types.Session) {

	if info, err := GetInfo(configData, sessionData); err == nil {

		sessionData.MaxQuantity = functions.StrToFloat64(info.MaxQuantity)
		sessionData.MinQuantity = functions.StrToFloat64(info.MinQuantity)
		sessionData.StepSize = functions.StrToFloat64(info.StepSize)

		return

	}

}

// GetSymbolFiatFunds Retrieve symbol fiat funds available
func GetSymbolFiatFunds(
	configData *types.Config,
	sessionData *types.Session) (balance float64, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceGetSymbolFiatFunds(sessionData)

	}

	return

}

// GetSymbolFunds Retrieve symbol funds available
func GetSymbolFunds(
	configData *types.Config,
	sessionData *types.Session) (balance float64, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceGetSymbolFunds(sessionData)

	}

	return

}

// GetKlines Retrieve KLines via REST API
func GetKlines(
	configData *types.Config,
	sessionData *types.Session) (klines []*types.Kline, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		tmp, err := binanceGetKlines(sessionData)

		if err == nil {
			return binanceMapKline(tmp), err
		}

		return nil, err

	}

	return

}

// GetPriceChangeStats Retrieve 24hs Rolling Price Statistics
func GetPriceChangeStats(
	configData *types.Config,
	sessionData *types.Session,
	marketData *types.Market) (priceChangeStats []*types.PriceChangeStats, err error) {

	switch strings.ToLower(configData.ExchangeName) {
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

/* Calculate the correct quantity to BUY according to the exchange lotSizeStep */
func getBuyQuantity(
	marketData *types.Market,
	sessionData *types.Session,
	fiatQuantity float64) (quantity float64) {

	return math.Round((fiatQuantity/marketData.Price)/sessionData.StepSize) * sessionData.StepSize

}

// GetUserStreamServiceListenKey Retrieve listen key for user stream service
func GetUserStreamServiceListenKey(
	configData *types.Config,
	sessionData *types.Session) (listenKey string, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceGetUserStreamServiceListenKey(sessionData)

	}

	return

}

// KeepAliveUserStreamServiceListenKey Keep user stream service alive
func KeepAliveUserStreamServiceListenKey(
	configData *types.Config,
	sessionData *types.Session) (err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceKeepAliveUserStreamServiceListenKey(sessionData)

	}

	return

}

// NewSetServerTimeService Synchronize time
func NewSetServerTimeService(
	configData *types.Config,
	sessionData *types.Session) (err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceNewSetServerTimeService(sessionData)

	}

	return

}

// WsBookTickerServe serve websocket that pushes updates to the best bid or ask price or quantity in real-time for a specified symbol.
func WsBookTickerServe(
	configData *types.Config,
	sessionData *types.Session,
	wsHandler *types.WsHandler,
	errHandler func(err error)) (doneC chan struct{}, stopC chan struct{}, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceWsBookTickerServe(sessionData, wsHandler, errHandler)

	}

	return

}

// WsKlineServe serve websocket kline handler
func WsKlineServe(
	configData *types.Config,
	sessionData *types.Session,
	wsHandler *types.WsHandler,
	errHandler func(err error)) (doneC chan struct{}, stopC chan struct{}, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceWsKlineServe(sessionData, wsHandler, errHandler)

	}

	return

}

// WsUserDataServe serve user data handler with listen key
func WsUserDataServe(
	configData *types.Config,
	sessionData *types.Session,
	wsHandler *types.WsHandler,
	errHandler func(err error)) (doneC chan struct{}, stopC chan struct{}, err error) {

	switch strings.ToLower(configData.ExchangeName) {
	case "binance":

		return binanceWsUserDataServe(sessionData, wsHandler, errHandler)

	}

	return

}

// BuyTicker Buy Ticker
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
	if configData.DryRun {

		logger.LogEntry{ /* Log Entry */
			Config:  configData,
			Market:  marketData,
			Session: sessionData,
			Order: &types.Order{
				Price: marketData.Price,
			},
			Message:  "BUYDRYRUN",
			LogLevel: "InfoLevel",
		}.Do()

		return

	}

	orderResponse, err := BuyOrder(
		configData,
		sessionData,
		functions.Float64ToStr(getBuyQuantity(marketData, sessionData, quantity), 4)) /* Get the correct quantity according to lotSizeMin and lotSizeStep */

	/* Test orderResponse for  errors */
	if (orderResponse == nil && err != nil) ||
		(orderResponse == nil && err == nil) {

		switch {
		case strings.Contains(err.Error(), "1013"):
			/* <APIError> code=-1013, msg=Filter failure: LOT_SIZE */

			/* Retrieve exchange lot size for ticker and store in sessionData */
			GetLotSize(configData, sessionData)

			return

		}

		return

	}

	/* Check if result is nil and set as zero */
	if orderPrice = orderResponse.CumulativeQuoteQuantity / orderResponse.ExecutedQuantity; math.IsNaN(orderPrice) {
		orderPrice = 0
	}

	orderExecutedQuantity = orderResponse.ExecutedQuantity

	/* Save order to database */
	if err := mysql.SaveOrder(
		sessionData,
		orderResponse,
		0, /* OrderIDSource */
		orderPrice /* OrderPrice */); err != nil {

		/* Cleanly exit ThreadID */
		threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

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

			orderPrice = orderStatus.CumulativeQuoteQuantity / orderStatus.ExecutedQuantity

			orderExecutedQuantity = orderStatus.ExecutedQuantity

			/* Update order status and price & Save Thread Transaction */
			if err := mysql.UpdateOrder(
				sessionData,
				int64(orderResponse.OrderID),
				orderResponse.CumulativeQuoteQuantity,
				orderResponse.ExecutedQuantity,
				orderPrice,
				string(orderStatus.Status)); err != nil {

				/* Cleanly exit ThreadID */
				threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

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
			orderResponse.CumulativeQuoteQuantity,
			orderPrice,
			orderExecutedQuantity); err != nil {

			/* Cleanly exit ThreadID */
			threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

		}

		logger.LogEntry{ /* Log Entry */
			Config:  configData,
			Market:  marketData,
			Session: sessionData,
			Order: &types.Order{
				OrderID: int(orderResponse.OrderID),
				Price:   orderPrice,
			},
			Message:  "BUY",
			LogLevel: "InfoLevel",
		}.Do()

	} else if isCanceled {

		logger.LogEntry{ /* Log Entry */
			Config:  configData,
			Market:  marketData,
			Session: sessionData,
			Order: &types.Order{
				OrderID: int(orderResponse.OrderID),
				Price:   orderPrice,
			},
			Message:  "CANCELED",
			LogLevel: "InfoLevel",
		}.Do()

	}

}

// SellTicker Sell Ticker
func SellTicker(
	order types.Order,
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) {

	var orderResponse *types.Order
	var orderStatus *types.Order

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
	if configData.DryRun {

		logger.LogEntry{ /* Log Entry */
			Config:  configData,
			Market:  marketData,
			Session: sessionData,
			Order: &types.Order{
				Price: marketData.Price,
			},
			Message:  "SELLDRYRUN",
			LogLevel: "InfoLevel",
		}.Do()

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

		logger.LogEntry{ /* Log Entry */
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		return

	}

	/* Save order to database */
	if err := mysql.SaveOrder(
		sessionData,
		orderResponse,
		int64(order.OrderID), /* OrderIDSource */
		marketData.Price /* OrderPrice */); err != nil {

		/* Cleanly exit ThreadID */
		threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

	}

S:
	switch orderResponse.Status {
	case "FILLED":

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
				threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

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
					case strings.Contains(err.Error(), "-2010"), strings.Contains(err.Error(), "-2011"), strings.Contains(err.Error(), "-1021"):
						/* -2011 Order filled in full before cancelling */
						/* -2010 Account has insufficient balance for requested action */
						/* -1021 Timestamp for this request was 1000ms ahead of the server's time */

						if orderStatus, err = GetOrder(
							configData,
							sessionData,
							int64(orderResponse.OrderID)); err != nil {

							/* Cleanly exit ThreadID */
							threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

						}

						break F

					case strings.Contains(err.Error(), "connection reset by peer"):
						/* read tcp 192.168.110.110:54914->65.9.137.130:443: read: connection reset by peer */

						if orderStatus, err = GetOrder(
							configData,
							sessionData,
							int64(orderResponse.OrderID)); err != nil {

							/* Cleanly exit ThreadID */
							threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

						}

						break S

					default:

						logger.LogEntry{ /* Log Entry */
							Config:  configData,
							Market:  marketData,
							Session: sessionData,
							Order: &types.Order{
								OrderID: int(orderResponse.OrderID),
							},
							Message:  functions.GetFunctionName() + " - " + err.Error(),
							LogLevel: "DebugLevel",
						}.Do()

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
						threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

					}

					break F

				default:

					logger.LogEntry{ /* Log Entry */
						Config:  configData,
						Market:  marketData,
						Session: sessionData,
						Order: &types.Order{
							OrderID: int(orderResponse.OrderID),
							Price:   marketData.Price,
						},
						Message:  "FAILED TO CANCEL ORDER",
						LogLevel: "InfoLevel",
					}.Do()

					break F

				}

			}

			/* Wait time between iterations (i++). There are ten iterations and the total waiting time define the amount od time before an order is canceled. configData.SellWaitBeforeCancel is divided by then converted into seconds. */
			time.Sleep(
				time.Duration(
					configData.SellWaitBeforeCancel/10) * time.Second)

		}

		/* Update order status and price */
		if err := mysql.UpdateOrder(
			sessionData,
			int64(orderResponse.OrderID),
			orderStatus.CumulativeQuoteQuantity,
			orderStatus.ExecutedQuantity,
			marketData.Price,
			string(orderStatus.Status)); err != nil {

			/* Cleanly exit ThreadID */
			threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

		}

	}

	if !isCanceled {

		/* Remove Thread transaction from database */
		if err := mysql.DeleteThreadTransactionByOrderID(
			sessionData,
			order.OrderID); err != nil {

			/* Cleanly exit ThreadID */
			threads.Thread{}.Terminate(sessionData, functions.GetFunctionName()+" - "+err.Error())

		}

		logger.LogEntry{ /* Log Entry */
			Config:  configData,
			Market:  marketData,
			Session: sessionData,
			Order: &types.Order{
				OrderID:       int(orderResponse.OrderID),
				Price:         marketData.Price,
				OrderIDSource: order.OrderID,
			},
			Message:  "SELL",
			LogLevel: "InfoLevel",
		}.Do()

	} else if isCanceled {

		logger.LogEntry{ /* Log Entry */
			Config:  configData,
			Market:  marketData,
			Session: sessionData,
			Order: &types.Order{
				OrderID:       int(orderResponse.OrderID),
				Price:         marketData.Price,
				OrderIDSource: order.OrderID,
			},
			Message:  "CANCELED",
			LogLevel: "InfoLevel",
		}.Do()

	}

}
