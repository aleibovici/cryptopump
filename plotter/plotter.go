package plotter

import (
	"bytes"
	"cryptopump/functions"
	"cryptopump/types"
	"html/template"
	"log"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	chartrender "github.com/go-echarts/go-echarts/v2/render"
)

// LoadKlineData load Kline into long term retention for plotter
func LoadKlineData(
	sessionData *types.Session,
	kline types.WsKline) {

	kd := []types.KlineData{
		{
			Date:    kline.EndTime,
			Data:    [4]float64{functions.StrToFloat64(kline.Open), functions.StrToFloat64(kline.Close), functions.StrToFloat64(kline.Low), functions.StrToFloat64(kline.High)},
			Volumes: functions.StrToFloat64(kline.Volume)},
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
		log.Printf("Failed to render chart: %s", err)
		return ""
	}

	return template.HTML(buf.String())
}

// Plot is responsible for rending e-chart
func Plot(sessionData *types.Session) (
	htmlSnippet template.HTML) {

	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	// v := make([]opts.BarData, 0)

	for i := 0; i < len(sessionData.KlineData); i++ {
		x = append(x, time.Unix((sessionData.KlineData[i].Date/1000), 0).UTC().Local().Format("15:04"))
		y = append(y, opts.KlineData{Value: sessionData.KlineData[i].Data})
		// v = append(v, opts.BarData{Value: sessionData.KlineData[i].Volumes})
	}

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
			Start:      60,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      60,
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

	kline.SetXAxis(x).AddSeries("kline", y).
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

	return renderToHTML(kline)
}
