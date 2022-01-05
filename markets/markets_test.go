package markets

import (
	"testing"

	"github.com/aleibovici/cryptopump/exchange"
	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/types"
	"github.com/sdcoffey/techan"
	"github.com/spf13/viper"
)

var configData = &types.Config{}

var sessionData = &types.Session{}

var marketData = &types.Market{
	Series: &techan.TimeSeries{},
}

func init() {

	viperData := &types.ViperData{ /* Viper Configuration */
		V1: viper.New(), /* Session configurations file */
		V2: viper.New(), /* Global configurations file */
	}

	viperData.V1.SetConfigType("yml")       /* Set the type of the configurations file */
	viperData.V1.AddConfigPath("../config") /* Set the path to look for the configurations file */
	viperData.V1.SetConfigName("config")    /* Set the file name of the configurations file */
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

	configData = functions.GetConfigData(viperData, sessionData)
	configData.TestNet = true

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

	sessionData.Symbol = configData.Symbol

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Data{
				Kline: tt.fields.Kline,
			}
			d.LoadKline(tt.args.configData, tt.args.sessionData, tt.args.marketData)
		})
	}
}
