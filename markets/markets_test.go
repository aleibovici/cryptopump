package markets

import (
	"database/sql"
	"testing"
	"time"

	"github.com/aleibovici/cryptopump/exchange"
	"github.com/aleibovici/cryptopump/types"
	"github.com/sdcoffey/techan"
)

var configData = &types.Config{
	ApikeyTestNet:    "i8kImZxWu9sgptrzatanbcoG2lgneGooDVoqSoNHjS3cdcEySqe0nG4NxNZ0WP4O",
	SecretkeyTestNet: "hglipvMt5t5NGu6aqwt1dgeTz5ss9cZUTO82IXo0thpBs8uBTmCea4xlJNIXdVgf",
	Symbol:           "BTCUSDT",
	ExchangeName:     "binance",
	TestNet:          true,
}

var sessionData = &types.Session{
	ThreadCount:          0,
	SellTransactionCount: 0,
	Symbol:               "BTCUSDT",
	SymbolFunds:          0,
	SymbolFiat:           "USDT",
	SymbolFiatFunds:      0,
	LastBuyTransactTime:  time.Time{},
	LastSellCanceledTime: time.Time{},
	ConfigTemplate:       0,
	ForceBuy:             false,
	ForceSell:            false,
	ListenKey:            "",
	MasterNode:           false,
	Db:                   &sql.DB{},
	Clients:              types.Client{},
	KlineData:            []types.KlineData{},
	StopWs:               false,
	Busy:                 false,
	MinQuantity:          0,
	MaxQuantity:          0,
	StepSize:             0,
}

var marketData = &types.Market{
	Rsi3:                      0,
	Rsi7:                      0,
	Rsi14:                     0,
	MACD:                      0,
	Price:                     40000,
	PriceChangeStatsHighPrice: 0,
	PriceChangeStatsLowPrice:  0,
	Direction:                 0,
	TimeStamp:                 time.Time{},
	Series:                    &techan.TimeSeries{},
	Ma7:                       0,
	Ma14:                      0,
}

func init() {

	exchange.GetClient(configData, sessionData)

}

func TestData_LoadKlinePast(t *testing.T) {
	type fields struct {
		Kline types.WsKline
	}
	type args struct {
		configData  *types.Config
		marketData  *types.Market
		sessionData *types.Session
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				Kline: types.WsKline{},
			},
			args: args{
				configData:  configData,
				marketData:  marketData,
				sessionData: sessionData,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Data{
				Kline: tt.fields.Kline,
			}
			d.LoadKlinePast(tt.args.configData, tt.args.marketData, tt.args.sessionData)
		})
	}
}

func TestData_LoadKline(t *testing.T) {
	type fields struct {
		Kline types.WsKline
	}
	type args struct {
		configData  *types.Config
		sessionData *types.Session
		marketData  *types.Market
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				Kline: types.WsKline{},
			},
			args: args{
				configData:  configData,
				marketData:  marketData,
				sessionData: sessionData,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Data{
				Kline: tt.fields.Kline,
			}
			d.LoadKline(tt.args.configData, tt.args.sessionData, tt.args.marketData)
		})
	}
}
