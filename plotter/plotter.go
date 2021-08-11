package plotter

import (
	"bytes"
	"html/template"
	"math"
	"time"

	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/types"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	chartrender "github.com/go-echarts/go-echarts/v2/render"
)

// LoadKlineData load Kline into long term retention for plotter
func LoadKlineData(
	sessionData *types.Session,
	marketData *types.Market,
	kline types.WsKline) {

	kd := []types.KlineData{
		{
			Date:    kline.EndTime,
			Data:    [4]float64{functions.StrToFloat64(kline.Open), functions.StrToFloat64(kline.Close), functions.StrToFloat64(kline.Low), functions.StrToFloat64(kline.High)},
			Volumes: functions.StrToFloat64(kline.Volume),
			Ma7:     math.Round(marketData.Ma7*10000) / 10000,
			Ma14:    math.Round(marketData.Ma14*10000) / 10000,
		},
	}

	/* Maintain klinedata to a maximum of 1440 minutes (24hs) eliminating first item on slice */
	if len(sessionData.KlineData) == 1439 {

		sessionData.KlineData = sessionData.KlineData[1:]

	}

	sessionData.KlineData = append(sessionData.KlineData, kd...)

}

func renderToHTML(c interface{}) template.HTML {
	var buf bytes.Buffer
	r := c.(chartrender.Renderer)
	err := r.Render(&buf)
	if err != nil {

		logger.LogEntry{
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		return ""
	}

	return template.HTML(buf.String())
}

// Plot is responsible for rending e-chart
func Plot(sessionData *types.Session) (
	htmlSnippet template.HTML) {

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	v := make([]opts.BarData, 0)
	ma7 := make([]opts.LineData, 0)  /* Simple Moving Average for 7 periods */
	ma14 := make([]opts.LineData, 0) /* Simple Moving Average for 14 periods */

	for i := 0; i < len(sessionData.KlineData); i++ {
		x = append(x, time.Unix((sessionData.KlineData[i].Date/1000), 0).UTC().Local().Format("15:04"))
		y = append(y, opts.KlineData{Value: sessionData.KlineData[i].Data})
		v = append(v, opts.BarData{Value: sessionData.KlineData[i].Volumes})
		ma7 = append(ma7, opts.LineData{Value: sessionData.KlineData[i].Ma7})    /* Simple Moving Average for 7 periods */
		ma14 = append(ma14, opts.LineData{Value: sessionData.KlineData[i].Ma14}) /* Simple Moving Average for 14 periods */
	}

	kline := klineBase("KLINE", x, y)                                                   /* Create base kline chart */
	kline.Overlap(lineBase("MA7", x, ma7, "blue"), lineBase("MA14", x, ma14, "orange")) /* Create overlaping line charts */

	return renderToHTML(kline)
}

/* Create a base KLine Chart */
func klineBase(name string, XAxis []string, klineData []opts.KlineData) *charts.Kline {

	kline := charts.NewKLine()

	kline.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "",
		}),
		charts.WithXAxisOpts(opts.XAxis{ /* Dates */
			Type:        "category",
			Show:        true,
			SplitNumber: 20,
			SplitArea:   &opts.SplitArea{},
			SplitLine:   &opts.SplitLine{},
			AxisLabel: &opts.AxisLabel{
				Show: true,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{ /* Candles */
			Type:        "value",
			SplitNumber: 2,
			Scale:       true,
			SplitArea: &opts.SplitArea{
				Show:      true,
				AreaStyle: &opts.AreaStyle{},
			},
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Color: "#777",
				},
			},
			AxisLabel: &opts.AxisLabel{
				Show:      true,
				Formatter: "{value}\n",
			},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      80,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      80,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: "CryptoPump",
			Width:     "1900px",
			Height:    "400px",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "axis",
			AxisPointer: &opts.AxisPointer{
				Type: "cross",
			},
		}),
	)

	kline.SetXAxis(XAxis).AddSeries(name, klineData).
		SetSeriesOptions(
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "highest value",
				Type:     "max",
				ValueDim: "highest",
			}),
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "lowest value",
				Type:     "min",
				ValueDim: "lowest",
			}),
			charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
				Label: &opts.Label{
					Show:  true,
					Color: "black",
				},
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:  "#00da3c",
				Color0: "#ec0000",
			}),
		)

	return kline

}

/* Create a base Line Chart */
func lineBase(name string, XAxis []string, lineData []opts.LineData, color string) *charts.Line {

	line := charts.NewLine()
	line.SetXAxis(XAxis).
		AddSeries(name, lineData).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: false,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   color,
				Width:   2,
				Opacity: 0.5,
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:        color,
				Color0:       color,
				BorderColor:  color,
				BorderColor0: color,
				Opacity:      0.5,
			}),
		)

	return line
}
