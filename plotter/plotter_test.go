package plotter

import (
	"html/template"
	"testing"

	"github.com/aleibovici/cryptopump/types"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func Test_klineBase(t *testing.T) {
	type args struct {
		name      string
		XAxis     []string
		klineData []opts.KlineData
	}
	tests := []struct {
		name string
		args args
		want *charts.Kline
	}{
		{
			name: "success",
			args: args{
				name:      "",
				XAxis:     []string{},
				klineData: []opts.KlineData{},
			},
			want: &charts.Kline{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := klineBase(tt.args.name, tt.args.XAxis, tt.args.klineData); got != nil {
				return
			}
		})
	}
}

func Test_lineBase(t *testing.T) {
	type args struct {
		name     string
		XAxis    []string
		lineData []opts.LineData
		color    string
	}
	tests := []struct {
		name string
		args args
		want *charts.Line
	}{
		{
			name: "success",
			args: args{
				name:     "MA7",
				XAxis:    make([]string, 0),
				lineData: make([]opts.LineData, 0),
				color:    "blue",
			},
			want: &charts.Line{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lineBase(tt.args.name, tt.args.XAxis, tt.args.lineData, tt.args.color); got != nil {
				return
			}
		})
	}
}

func TestData_Plot(t *testing.T) {
	type fields struct {
		Kline types.WsKline
	}
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantHTMLSnippet template.HTML
	}{
		{
			name: "success",
			fields: fields{
				Kline: types.WsKline{},
			},
			args: args{
				sessionData: &types.Session{
					KlineData: []types.KlineData{},
				},
			},
			wantHTMLSnippet: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Data{
				Kline: tt.fields.Kline,
			}
			if gotHTMLSnippet := d.Plot(tt.args.sessionData); gotHTMLSnippet != "" {
				return
			}
		})
	}
}

func TestData_LoadKline(t *testing.T) {
	type fields struct {
		Kline types.WsKline
	}
	type args struct {
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
				Kline: types.WsKline{
					StartTime:            1630474980000,
					EndTime:              1630475039999,
					Symbol:               "BTCUSDT",
					Interval:             "1m",
					FirstTradeID:         1291659,
					LastTradeID:          1291679,
					Open:                 "47397.49000000",
					Close:                "47372.73000000",
					High:                 "47397.49000000",
					Low:                  "47354.42000000",
					Volume:               "0.20152400",
					TradeNum:             21,
					IsFinal:              true,
					QuoteVolume:          "9547.12114436",
					ActiveBuyVolume:      "9499.76666436",
					ActiveBuyQuoteVolume: "9499.76666436",
				},
			},
			args: args{
				sessionData: &types.Session{},
				marketData:  &types.Market{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Data{
				Kline: tt.fields.Kline,
			}
			d.LoadKline(tt.args.sessionData, tt.args.marketData)
		})
	}
}
