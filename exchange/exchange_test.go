package exchange

import (
	"reflect"
	"testing"

	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/types"
	"github.com/spf13/viper"
)

var configData = &types.Config{}

var sessionData = &types.Session{
	Symbol:     "BTCUSDT",
	SymbolFiat: "USDT",
	MasterNode: false,
}

var marketData = &types.Market{
	Price: 40000,
}

func init() {

	viper.AddConfigPath("../config") /* Set the path to look for the configurations file */
	if err := viper.ReadInConfig(); err != nil {

		logger.LogEntry{
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	configData = functions.GetConfigData(sessionData)
	configData.TestNet = true

	GetClient(configData, sessionData)
	GetLotSize(configData, sessionData)

}

func TestGetClient(t *testing.T) {
	type args struct {
		configData  *types.Config
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				configData:  configData,
				sessionData: sessionData,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetClient(tt.args.configData, tt.args.sessionData); (err != nil) != tt.wantErr {
				t.Errorf("GetClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetInfo(t *testing.T) {
	type args struct {
		configData  *types.Config
		sessionData *types.Session
	}

	tests := []struct {
		name     string
		args     args
		wantInfo *types.ExchangeInfo
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				configData:  configData,
				sessionData: sessionData},
			wantInfo: &types.ExchangeInfo{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, err := GetInfo(tt.args.configData, tt.args.sessionData)
			if err == nil {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("GetInfo() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func TestGetLotSize(t *testing.T) {
	type args struct {
		configData  *types.Config
		sessionData *types.Session
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				configData:  configData,
				sessionData: sessionData,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetLotSize(tt.args.configData, tt.args.sessionData)
		})
	}
}

func TestGetKlines(t *testing.T) {
	type args struct {
		configData  *types.Config
		sessionData *types.Session
	}

	tests := []struct {
		name       string
		args       args
		wantKlines []*types.Kline
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				configData:  configData,
				sessionData: sessionData,
			},
			wantKlines: []*types.Kline{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKlines, err := GetKlines(tt.args.configData, tt.args.sessionData)
			if err == nil {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKlines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKlines, tt.wantKlines) {
				t.Errorf("GetKlines() = %v, want %v", gotKlines, tt.wantKlines)
			}
		})
	}
}

func TestGetPriceChangeStats(t *testing.T) {
	type args struct {
		configData  *types.Config
		sessionData *types.Session
		marketData  *types.Market
	}

	tests := []struct {
		name                 string
		args                 args
		wantPriceChangeStats []*types.PriceChangeStats
		wantErr              bool
	}{
		{
			name: "success",
			args: args{
				configData:  configData,
				sessionData: sessionData,
				marketData:  marketData,
			},
			wantPriceChangeStats: []*types.PriceChangeStats{},
			wantErr:              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPriceChangeStats, err := GetPriceChangeStats(tt.args.configData, tt.args.sessionData, tt.args.marketData)
			if err == nil {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPriceChangeStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPriceChangeStats, tt.wantPriceChangeStats) {
				t.Errorf("GetPriceChangeStats() = %v, want %v", gotPriceChangeStats, tt.wantPriceChangeStats)
			}
		})
	}
}

func TestGetUserStreamServiceListenKey(t *testing.T) {
	type args struct {
		configData  *types.Config
		sessionData *types.Session
	}

	tests := []struct {
		name          string
		args          args
		wantListenKey string
		wantErr       bool
	}{
		{
			name: "success",
			args: args{
				configData:  configData,
				sessionData: sessionData,
			},
			wantListenKey: "",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotListenKey, err := GetUserStreamServiceListenKey(tt.args.configData, tt.args.sessionData)
			if (err != nil) && (gotListenKey != tt.wantListenKey) {
				return
			}
		})
	}
}

func TestKeepAliveUserStreamServiceListenKey(t *testing.T) {
	type args struct {
		configData  *types.Config
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				configData:  configData,
				sessionData: sessionData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := KeepAliveUserStreamServiceListenKey(tt.args.configData, tt.args.sessionData); (err != nil) != tt.wantErr {
				t.Errorf("KeepAliveUserStreamServiceListenKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewSetServerTimeService(t *testing.T) {
	type args struct {
		configData  *types.Config
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				configData:  configData,
				sessionData: sessionData,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewSetServerTimeService(tt.args.configData, tt.args.sessionData); (err != nil) != tt.wantErr {
				t.Errorf("NewSetServerTimeService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getSellQuantity(t *testing.T) {
	type args struct {
		order       types.Order
		sessionData *types.Session
	}
	tests := []struct {
		name         string
		args         args
		wantQuantity float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotQuantity := getSellQuantity(tt.args.order, tt.args.sessionData); gotQuantity != tt.wantQuantity {
				t.Errorf("getSellQuantity() = %v, want %v", gotQuantity, tt.wantQuantity)
			}
		})
	}
}

func Test_getBuyQuantity(t *testing.T) {
	type args struct {
		marketData   *types.Market
		sessionData  *types.Session
		fiatQuantity float64
	}
	tests := []struct {
		name         string
		args         args
		wantQuantity float64
	}{
		{
			name: "success",
			args: args{
				marketData:   marketData,
				sessionData:  sessionData,
				fiatQuantity: 0,
			},
			wantQuantity: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotQuantity := getBuyQuantity(tt.args.marketData, tt.args.sessionData, tt.args.fiatQuantity); gotQuantity != tt.wantQuantity {
				t.Errorf("getBuyQuantity() = %v, want %v", gotQuantity, tt.wantQuantity)
			}
		})
	}
}
