package markets

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aleibovici/cryptopump/exchange"
	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/types"

	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

/* Technical analysis Calculations */
func calculate(
	closePrices techan.Indicator,
	priceChangeStats []*types.PriceChangeStats,
	sessionData *types.Session,
	marketData *types.Market) {

	marketData.Rsi3 = calculateRSI(closePrices, marketData.Series, 3)
	marketData.Rsi7 = calculateRSI(closePrices, marketData.Series, 7)
	marketData.Rsi14 = calculateRSI(closePrices, marketData.Series, 14)
	marketData.MACD = calculateMACD(closePrices, marketData.Series, 12, 26)
	marketData.Ma7 = calculateMA(closePrices, marketData.Series, 7)
	marketData.Ma14 = calculateMA(closePrices, marketData.Series, 14)
	if priceChangeStats != nil {
		marketData.PriceChangeStatsHighPrice = calculatePriceChangeStatsHighPrice(priceChangeStats)
		marketData.PriceChangeStatsLowPrice = calculatePriceChangeStatsLowPrice(priceChangeStats)
	}
	marketData.TimeStamp = time.Now() /* Time of last retrieved market Data */

}

// LoadKlineData Retrieve RealTime Kline Data
func LoadKlineData(
	configData *types.Config,
	sessionData *types.Session,
	marketData *types.Market,
	kline types.WsKline) {

	var start int64
	var err error
	var priceChangeStats []*types.PriceChangeStats

	if start, err = strconv.ParseInt(fmt.Sprint(kline.StartTime), 10, 64); err != nil {

		logger.LogEntry{
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		return

	}

	period := techan.NewTimePeriod(time.Unix((start/1000), 0).UTC(), time.Minute*1)

	candle := techan.NewCandle(period)
	candle.OpenPrice = big.NewFromString(kline.Open)
	candle.ClosePrice = big.NewFromString(kline.Close)
	candle.MaxPrice = big.NewFromString(kline.High)
	candle.MinPrice = big.NewFromString(kline.Low)
	candle.Volume = big.NewFromString(kline.Volume)

	if !marketData.Series.AddCandle(candle) { /* AddCandle adds the given candle to TimeSeries */

		return

	}

	if priceChangeStats, err = exchange.GetPriceChangeStats(configData, sessionData, marketData); err != nil {

		logger.LogEntry{
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		return

	}

	calculate(
		techan.NewClosePriceIndicator(marketData.Series),
		priceChangeStats,
		sessionData,
		marketData)

}

// LoadKlineDataPast Retrieve Old Kline Data
func LoadKlineDataPast(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session) {

	var err error
	var klines []*types.Kline
	var priceChangeStats []*types.PriceChangeStats

	if klines, err = exchange.GetKlines(configData, sessionData); err != nil {

		return

	}

	for _, datum := range klines {

		var start int64

		if start, err = strconv.ParseInt(fmt.Sprint(datum.OpenTime), 10, 64); err != nil {

			logger.LogEntry{
				Config:   configData,
				Market:   marketData,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "DebugLevel",
			}.Do()

			return

		}

		period := techan.NewTimePeriod(time.Unix((start/1000), 0).UTC(), time.Minute*1)

		candle := techan.NewCandle(period)
		candle.OpenPrice = big.NewFromString(datum.Open)
		candle.ClosePrice = big.NewFromString(datum.Close)
		candle.MaxPrice = big.NewFromString(datum.High)
		candle.MinPrice = big.NewFromString(datum.Low)
		candle.Volume = big.NewFromString(datum.Volume)

		if !marketData.Series.AddCandle(candle) {
			return
		}

	}

	if priceChangeStats, err = exchange.GetPriceChangeStats(configData, sessionData, marketData); err != nil {

		logger.LogEntry{
			Config:   configData,
			Market:   marketData,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		return

	}

	calculate(
		techan.NewClosePriceIndicator(marketData.Series),
		priceChangeStats,
		sessionData,
		marketData)

}

/* Calculate Relative Strength Index */
func calculateRSI(
	closePrices techan.Indicator,
	series *techan.TimeSeries,
	timeframe int) float64 {

	return techan.NewRelativeStrengthIndexIndicator(closePrices, timeframe).Calculate(series.LastIndex() - 1).Float()
}

func calculateMACD(
	closePrices techan.Indicator,
	series *techan.TimeSeries,
	shortwindow int,
	longwindow int) float64 {

	return techan.NewMACDIndicator(closePrices, shortwindow, longwindow).Calculate(series.LastIndex() - 1).Float()
}

func calculateMA(
	closePrices techan.Indicator,
	series *techan.TimeSeries,
	window int) float64 {

	return techan.NewSimpleMovingAverage(closePrices, window).Calculate(series.LastIndex() - 1).Float()
}

/* Calculate High price for 1 period */
func calculatePriceChangeStatsHighPrice(
	priceChangeStats []*types.PriceChangeStats) float64 {

	return functions.StrToFloat64(priceChangeStats[0].HighPrice)
}

/* Calculate Low price for 1 period */
func calculatePriceChangeStatsLowPrice(
	priceChangeStats []*types.PriceChangeStats) float64 {

	return functions.StrToFloat64(priceChangeStats[0].LowPrice)
}
