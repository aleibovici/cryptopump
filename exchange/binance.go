package exchange

import (
	"context"
	"time"

	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/types"

	"github.com/adshao/go-binance/v2"
)

/* Map binance.Order types to Order type */
func binanceMapOrder(from *binance.Order) (to *types.Order) {

	to = &types.Order{}
	to.ClientOrderID = from.ClientOrderID
	to.OrderID = int(from.OrderID)
	to.CumulativeQuoteQuantity = functions.StrToFloat64(from.CummulativeQuoteQuantity)
	to.ExecutedQuantity = functions.StrToFloat64(from.ExecutedQuantity)
	to.Price = functions.StrToFloat64(from.Price)
	to.Side = string(from.Side)
	to.Status = string(from.Status)
	to.Symbol = from.Symbol

	return to

}

/* Map binance.CreateOrderResponse types to Order type */
func binanceMapCreateOrderResponse(from *binance.CreateOrderResponse) (to *types.Order) {

	to = &types.Order{}
	to.ClientOrderID = from.ClientOrderID
	to.OrderID = int(from.OrderID)
	to.CumulativeQuoteQuantity = functions.StrToFloat64(from.CummulativeQuoteQuantity)
	to.ExecutedQuantity = functions.StrToFloat64(from.ExecutedQuantity)
	to.Price = functions.StrToFloat64(from.Price)
	to.Side = string(from.Side)
	to.Status = string(from.Status)
	to.Symbol = from.Symbol
	to.TransactTime = from.TransactTime

	return to

}

/* Map binance.CancelOrderResponse types to Order type */
func binanceMapCancelOrderResponse(from *binance.CancelOrderResponse) (to *types.Order) {

	to = &types.Order{}
	to.ClientOrderID = from.ClientOrderID
	to.OrderID = int(from.OrderID)
	to.CumulativeQuoteQuantity = functions.StrToFloat64(from.CummulativeQuoteQuantity)
	to.ExecutedQuantity = functions.StrToFloat64(from.ExecutedQuantity)
	to.Price = functions.StrToFloat64(from.Price)
	to.Side = string(from.Side)
	to.Status = string(from.Status)
	to.Symbol = from.Symbol
	to.TransactTime = from.TransactTime

	return to

}

/* Map binance.Kline types to Kline type */
func binanceMapKline(from []*binance.Kline) (to []*types.Kline) {

	to = []*types.Kline{}

	for key := range from {

		tmp := &types.Kline{}
		tmp.Close = from[key].Close
		tmp.High = from[key].High
		tmp.Low = from[key].Low
		tmp.Open = from[key].Open
		tmp.OpenTime = from[key].OpenTime
		tmp.Volume = from[key].Volume

		to = append(to, tmp)

	}

	return to

}

// BinanceMapWsKline Map binance.WsKline types to WsKline type
func BinanceMapWsKline(from binance.WsKline) (to types.WsKline) {

	to = types.WsKline{}
	to.ActiveBuyQuoteVolume = from.ActiveBuyQuoteVolume
	to.ActiveBuyVolume = from.ActiveBuyQuoteVolume
	to.Close = from.Close
	to.EndTime = from.EndTime
	to.FirstTradeID = from.FirstTradeID
	to.High = from.High
	to.Interval = from.Interval
	to.IsFinal = from.IsFinal
	to.LastTradeID = from.LastTradeID
	to.Low = from.Low
	to.Open = from.Open
	to.QuoteVolume = from.QuoteVolume
	to.StartTime = from.StartTime
	to.Symbol = from.Symbol
	to.TradeNum = from.TradeNum
	to.Volume = from.Volume

	return to

}

/* Map binance.PriceChangeStats types to Kline type */
func binanceMapPriceChangeStats(from []*binance.PriceChangeStats) (to []*types.PriceChangeStats) {

	to = []*types.PriceChangeStats{}

	for key := range from {

		tmp := &types.PriceChangeStats{}
		tmp.HighPrice = from[key].HighPrice
		tmp.LowPrice = from[key].LowPrice

		to = append(to, tmp)

	}

	return to

}

/* Map binance.ExchangeInfo types to Order type */
func binanceMapExchangeInfo(sessionData *types.Session, from *binance.ExchangeInfo) (to *types.ExchangeInfo) {

	to = &types.ExchangeInfo{}

	for key := range from.Symbols {

		if from.Symbols[key].Symbol == sessionData.Symbol {

			to.MaxQuantity = from.Symbols[key].LotSizeFilter().MaxQuantity
			to.MinQuantity = from.Symbols[key].LotSizeFilter().MinQuantity
			to.StepSize = from.Symbols[key].LotSizeFilter().StepSize

		}

	}

	return to

}

/* Get Binance client */
func binanceGetClient(
	configData *types.Config) *binance.Client {

	binance.WebsocketKeepalive = false
	binance.WebsocketTimeout = time.Second * 30

	/* Exchange test network, used with launch.json */
	if configData.TestNet {

		binance.UseTestnet = true
		return binance.NewClient(configData.ApikeyTestNet, configData.SecretkeyTestNet)

	}

	return binance.NewClient(configData.Apikey, configData.Secretkey)

}

/* Retrieve exchange information */
func binanceGetInfo(
	sessionData *types.Session) (info *types.ExchangeInfo, err error) {

	var tmp *binance.ExchangeInfo

	if tmp, err = sessionData.Clients.Binance.NewExchangeInfoService().Do(context.Background()); err != nil {

		return nil, err

	}

	return binanceMapExchangeInfo(sessionData, tmp), err

}

/* Retrieve listen key for user stream service */
func binanceGetUserStreamServiceListenKey(
	sessionData *types.Session) (listenKey string, err error) {

	if listenKey, err = sessionData.Clients.Binance.NewStartUserStreamService().Do(context.Background()); err != nil {

		return "", err

	}

	return listenKey, err

}

/* Keep user stream service alive */
func binanceKeepAliveUserStreamServiceListenKey(
	sessionData *types.Session) (err error) {

	if err = sessionData.Clients.Binance.NewKeepaliveUserStreamService().ListenKey(sessionData.ListenKey).Do(context.Background()); err != nil {

		return err

	}

	return

}

/* Synchronize time */
func binanceNewSetServerTimeService(
	sessionData *types.Session) (err error) {

	if _, err = sessionData.Clients.Binance.NewSetServerTimeService().Do(context.Background()); err != nil {

		return err

	}

	return

}

/* Retrieve funds available */
func binanceGetSymbolFunds(
	sessionData *types.Session) (balance float64, err error) {

	var account *binance.Account

	if account, err = sessionData.Clients.Binance.NewGetAccountService().Do(context.Background()); err != nil {

		logger.LogEntry{
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		return 0, err

	}

	for key := range account.Balances {

		if account.Balances[key].Asset == sessionData.SymbolFiat {

			return functions.StrToFloat64(account.Balances[key].Free), err

		}

	}

	return 0, err

}

/* Minutely crypto currency open/close prices, high/low, trades and others */
func binanceGetKlines(
	sessionData *types.Session) (klines []*binance.Kline, err error) {

	if klines, err = sessionData.Clients.Binance.NewKlinesService().Symbol(sessionData.Symbol).
		Interval("1m").Limit(14).Do(context.Background()); err != nil {

		return nil, err

	}

	return klines, err

}

/* 24hr ticker price change statistics */
func binanceGetPriceChangeStats(
	sessionData *types.Session) (PriceChangeStats []*types.PriceChangeStats, err error) {

	var tmp []*binance.PriceChangeStats

	if tmp, err = sessionData.Clients.Binance.NewListPriceChangeStatsService().Symbol(sessionData.Symbol).Do(context.Background()); err != nil {

		return nil, err

	}

	return binanceMapPriceChangeStats(tmp), err

}

/* Retrieve Order Status */
func binanceGetOrder(
	sessionData *types.Session,
	orderID int64) (order *types.Order, err error) {

	var tmp *binance.Order

	if tmp, err = sessionData.Clients.Binance.NewGetOrderService().Symbol(sessionData.Symbol).OrderID(orderID).Do(context.Background()); err != nil {

		return nil, err

	}

	return binanceMapOrder(tmp), err

}

/* CANCEL an order */
func binanceCancelOrder(
	sessionData *types.Session,
	orderID int64) (cancelOrderResponse *types.Order, err error) {

	var tmp *binance.CancelOrderResponse

	if tmp, err = sessionData.Clients.Binance.NewCancelOrderService().Symbol(sessionData.Symbol).OrderID(orderID).Do(context.Background()); err != nil {

		return nil, err

	}

	return binanceMapCancelOrderResponse(tmp), err

}

/* Create order to BUY */
func binanceBuyOrder(
	sessionData *types.Session,
	quantity string) (order *types.Order, err error) {

	var tmp *binance.CreateOrderResponse

	/* Execute OrderTypeMarket */
	if tmp, err = sessionData.Clients.Binance.NewCreateOrderService().Symbol(sessionData.Symbol).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).
		Quantity(quantity).Do(context.Background()); err != nil {

		logger.LogEntry{
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "InfoLevel",
		}.Do()

		return nil, err

	}

	return binanceMapCreateOrderResponse(tmp), err

}

/* WsBookTickerServe serve websocket that pushes updates to the best bid or ask price or quantity in real-time for a specified symbol. */
func binanceWsBookTickerServe(
	sessionData *types.Session,
	wsHandler *types.WsHandler,
	errHandler func(err error)) (doneC chan struct{}, stopC chan struct{}, err error) {

	doneC, stopC, err = binance.WsBookTickerServe(sessionData.Symbol, wsHandler.BinanceWsBookTicker, errHandler)

	return doneC, stopC, err

}

/* WsKlineServe serve websocket kline handler */
func binanceWsKlineServe(
	sessionData *types.Session,
	wsHandler *types.WsHandler,
	errHandler func(err error)) (doneC chan struct{}, stopC chan struct{}, err error) {

	doneC, stopC, err = binance.WsKlineServe(sessionData.Symbol, "1m", wsHandler.BinanceWsKline, errHandler)

	return doneC, stopC, err

}

/* WsUserDataServe serve user data handler with listen key */
func binanceWsUserDataServe(
	sessionData *types.Session,
	wsHandler *types.WsHandler,
	errHandler func(err error)) (doneC chan struct{}, stopC chan struct{}, err error) {

	doneC, stopC, err = binance.WsUserDataServe(sessionData.ListenKey, wsHandler.BinanceWsUserDataServe, errHandler)

	return doneC, stopC, err
}

/* Create order to SELL */
func binanceSellOrder(
	marketData *types.Market,
	sessionData *types.Session,
	quantity string) (order *types.Order, err error) {

	var tmp *binance.CreateOrderResponse

	if !sessionData.ForceSell {

		/* Execute OrderTypeLimit */
		if tmp, err = sessionData.Clients.Binance.NewCreateOrderService().Symbol(sessionData.Symbol).Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).Quantity(quantity).Price(functions.Float64ToStr(marketData.Price, 2)).TimeInForce(binance.TimeInForceTypeGTC).Do(context.Background()); err != nil {

			return nil, err

		}

	}

	if sessionData.ForceSell {

		sessionData.ForceSell = false

		/* Execute OrderTypeMarket */
		if tmp, err = sessionData.Clients.Binance.NewCreateOrderService().Symbol(sessionData.Symbol).Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).Quantity(quantity).Do(context.Background()); err != nil {

			return nil, err

		}

	}

	return binanceMapCreateOrderResponse(tmp), err
}
