package functions

import (
	"fmt"
	"html/template"
	"io"
	"math"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"os"
	"strconv"

	"cryptopump/types"

	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// StrToFloat64 function
/* This public function convert string to float64 */
func StrToFloat64(value string) (r float64) {

	var err error

	if r, err = strconv.ParseFloat(value, 8); err != nil {

		log.Fatal(err)

	}

	return r
}

// Float64ToStr function
/* This public function convert float64 to string with variable precision */
func Float64ToStr(value float64, prec int) string {

	return strconv.FormatFloat(value, 'f', prec, 64)

}

// IntToFloat64
/* This public function convert Int to Float64 */
func IntToFloat64(value int) float64 {

	return float64(value)

}

// Profit
/* This public function calculate profit */
func GetProfitResult(buyPrice float64, sellPrice float64) float64 {

	return (sellPrice - buyPrice) / buyPrice

}

// StrToInt function
/* This public function convert string to int */
func StrToInt(value string) (r int) {

	var err error

	if r, err = strconv.Atoi(value); err != nil {

		fmt.Println(err)

	}

	return r

}

func Logger(
	configData *types.Config,
	marketData *types.Market,
	sessionData *types.Session,
	logLevel log.Level,
	orderIDThread int,
	orderID int64,
	orderPrice float64,
	profit float64,
	message string) {

	var err error
	var filename string
	var file *os.File

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		DisableSorting:  false,
	})

	// Only log the warning severity or above.
	log.SetLevel(logLevel)

	switch {
	case logLevel == log.InfoLevel:

		filename = "cryptopump.log"

	case logLevel == log.DebugLevel:

		filename = "cryptopump_debug.log"

	}

	// You could set this to any `io.Writer` such as a file
	if file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666); err != nil {

		log.Fatal(err)

	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(file)

	switch {
	case logLevel == log.InfoLevel:

		switch message {
		case "UP", "DOWN", "INIT":

			log.WithFields(log.Fields{
				"threadID":  sessionData.ThreadID,
				"rsi3":      fmt.Sprintf("%.2f", marketData.Rsi3),
				"rsi7":      fmt.Sprintf("%.2f", marketData.Rsi7),
				"rsi14":     fmt.Sprintf("%.2f", marketData.Rsi14),
				"MACD":      fmt.Sprintf("%.2f", marketData.MACD),
				"high":      marketData.PriceChangeStatsHighPrice,
				"direction": marketData.Direction,
			}).Info(message)

		case "BUY":

			log.WithFields(log.Fields{
				"threadID":   sessionData.ThreadID,
				"orderID":    orderID,
				"orderPrice": fmt.Sprintf("%.4f", orderPrice),
			}).Info(message)

		case "SELL":

			log.WithFields(log.Fields{
				"threadID":      sessionData.ThreadID,
				"orderIDThread": orderIDThread,
				"orderID":       orderID,
				"orderPrice":    fmt.Sprintf("%.4f", orderPrice),
				"profit":        fmt.Sprintf("%.4f", profit),
			}).Info(message)

		case "CANCELED":

			if configData.Debug.(string) == "true" {

				log.WithFields(log.Fields{
					"threadID":      sessionData.ThreadID,
					"orderIDThread": orderIDThread,
					"orderID":       orderID,
				}).Info(message)

			}

		default:

			log.WithFields(log.Fields{
				"threadID": sessionData.ThreadID,
			}).Info(message)

		}

	case logLevel == log.DebugLevel:

		log.WithFields(log.Fields{
			"threadID": sessionData.ThreadID,
			"orderID":  orderID,
		}).Debug(message)

	}

}

// MustGetenv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func MustGetenv(k string) string {

	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}

	return v

}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {

	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}

/* Return random thread ID */
func GetThreadID() string {

	return xid.New().String()

}

/* Convert Strign to Time */
func stringToTime(str string) (r time.Time) {

	var err error

	if r, err = time.Parse(time.Kitchen, str); err != nil {

		Logger(
			nil,
			nil,
			nil,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			GetFunctionName()+" - "+err.Error())

	}

	return r

}

/* Check if time is in a specific range */
func IsInTimeRange(startTimeString string, endTimeString string) bool {

	t := time.Now()
	timeNowString := t.Format(time.Kitchen)
	timeNow := stringToTime(timeNowString)
	start := stringToTime(startTimeString)
	end := stringToTime(endTimeString)

	if timeNow.Before(start) {

		return false

	}

	if timeNow.After(end) {

		return false

	}

	return true

}

/* Validate available funds to buy */
func IsFundsAvailable(
	configData *types.Config,
	sessionData *types.Session) bool {

	return (sessionData.Symbol_fiat_funds - StrToFloat64(configData.Symbol_fiat_stash.(string))) >= StrToFloat64(configData.Buy_quantity_fiat_down.(string))

}

/* Convert Fiat currency to Coin */
func ConvertFiatToCoin(fiatQty float64, ticker_price float64, lotSizeMin float64, lotSizeStep float64) float64 {

	return math.Round((fiatQty/ticker_price)/lotSizeStep) * lotSizeStep

}

// ExecuteTemplate function
/* This public function i responsible for executing any templates */
func ExecuteTemplate(
	wr io.Writer,
	name string,
	data interface{},
	sessionData *types.Session) {

	var tlp *template.Template
	var err error

	if tlp, err = template.ParseGlob("./templates/*"); err != nil {

		Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			GetFunctionName()+" - "+err.Error())

		os.Exit(1)

	}

	if err = tlp.ExecuteTemplate(wr, name, data); err != nil {

		Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			GetFunctionName()+" - "+err.Error())

		os.Exit(1)

	}

}

/* Retrieve current function name */
func GetFunctionName() string {

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return frame.Function

}

/* Create lock for threadID */
func LockThreadID(threadID string) bool {

	filename := threadID + ".lock"

	if _, err := os.Stat(filename); err == nil {

		return false

	} else if os.IsNotExist(err) {

		var file, err = os.Create(filename)
		if err != nil {
			return false
		}

		file.Close()

		return true

	}

	return false

}

/* Determine port for HTTP service. */
func GetPort() (port string) {

	port = os.Getenv("PORT")

	if port == "" {

		port = "8080"

	}

	for {

		if l, err := net.Listen("tcp", ":"+port); err != nil {

			port = Float64ToStr((StrToFloat64(port) + 1), 0)

		} else {

			l.Close()
			break

		}

	}

	return port

}

/* Delete configuration file for ThreadID */
func DeleteConfigFile(sessionData *types.Session) {

	filename := sessionData.ThreadID + ".yml"
	path := "./config/"

	if err := os.Remove(path + filename); err != nil {

		Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			GetFunctionName()+" - "+err.Error())

		return

	}

}

/* Retrieve or create config file based on ThreadID */
func GetConfigData(
	sessionData *types.Session) *types.Config {

	configData := loadConfigData(sessionData)

	if sessionData.ThreadID != "" {

		filename := sessionData.ThreadID + ".yml"
		writePath := "./config/"

		if _, err := os.Stat(writePath + filename); err == nil {

			/* Test for existing ThreadID config file and load configuration */
			viper.SetConfigFile(writePath + filename)

			if err := viper.ReadInConfig(); err != nil {

				Logger(
					nil,
					nil,
					sessionData,
					log.DebugLevel,
					0,
					0,
					0,
					0,
					GetFunctionName()+" - "+err.Error())

			}

			configData = loadConfigData(sessionData)

		} else if os.IsNotExist(err) {

			/* Create new ThreadID config file and load configuration */
			if err := viper.WriteConfigAs(writePath + filename); err != nil {

				Logger(
					nil,
					nil,
					sessionData,
					log.DebugLevel,
					0,
					0,
					0,
					0,
					GetFunctionName()+" - "+err.Error())

			}

			viper.SetConfigFile(writePath + filename)

			if err := viper.ReadInConfig(); err != nil {

				Logger(
					nil,
					nil,
					sessionData,
					log.DebugLevel,
					0,
					0,
					0,
					0,
					GetFunctionName()+" - "+err.Error())

			}

			configData = loadConfigData(sessionData)

		}

	}

	return configData

}

/* This function retrieve the list of configuration files under the root config folder.
.yaml files are considered configuration files. */
func getConfigTemplateList(sessionData *types.Session) []string {

	var files []string
	files = append(files, "-")

	root := "./config"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if filepath.Ext(path) == ".yml" {
			files = append(files, info.Name())
		}

		return nil
	})

	if err != nil {

		Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			GetFunctionName()+" - "+err.Error())

		os.Exit(1)

	}

	return files

}

/* Load the selected configuration template */
/* Three's a BUG where it only works before the first UPDATE */
func LoadConfigTemplate(
	sessionData *types.Session) *types.Config {

	var filename string

	/* Retrieve the list of configuration templates */
	files := getConfigTemplateList(sessionData)

	/* Iterate configuration templates to match the selection */
	for key, file := range files {
		if key == sessionData.ConfigTemplate {
			filename = file
		}
	}

	filenameOld := viper.ConfigFileUsed()

	/* Set selected template as current config and load settings and return configData*/
	viper.SetConfigFile("./config/" + filename)
	if err := viper.ReadInConfig(); err != nil {

		Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			GetFunctionName()+" - "+err.Error())

	}

	configData := loadConfigData(sessionData)

	/* Set origina template as current config */
	viper.SetConfigFile(filenameOld)
	if err := viper.ReadInConfig(); err != nil {

		Logger(
			configData,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			GetFunctionName()+" - "+err.Error())

	}

	return configData

}

/* This routine load viper configuration data into map[string]interface{} */
func loadConfigData(
	sessionData *types.Session) *types.Config {

	configData := &types.Config{
		ThreadID:                         sessionData.ThreadID, /* For index.html population */
		Apikey:                           viper.GetString("config.apiKey"),
		Secretkey:                        viper.GetString("config.secretKey"),
		ApikeyTestNet:                    viper.GetString("config.apiKeyTestNet"),    /* API key for exchange test network, used with launch.json */
		SecretkeyTestNet:                 viper.GetString("config.secretKeyTestNet"), /* Secret key for exchange test network, used with launch.json */
		Buy_24hs_highprice_entry:         viper.GetString("config.buy_24hs_highprice_entry"),
		Buy_24hs_highprice_entry_MACD:    viper.GetString("config.buy_24hs_highprice_entry_MACD"),
		Buy_direction_down:               viper.GetString("config.buy_direction_down"),
		Buy_direction_up:                 viper.GetString("config.buy_direction_up"),
		Buy_quantity_fiat_up:             viper.GetString("config.buy_quantity_fiat_up"),
		Buy_quantity_fiat_down:           viper.GetString("config.buy_quantity_fiat_down"),
		Buy_quantity_fiat_init:           viper.GetString("config.buy_quantity_fiat_init"),
		Buy_repeat_threshold_down:        viper.GetString("config.buy_repeat_threshold_down"),
		Buy_repeat_threshold_down_second: viper.GetString("config.buy_repeat_threshold_down_second"),
		Buy_repeat_threshold_down_second_start_count: viper.GetString("config.buy_repeat_threshold_down_second_start_count"),
		Buy_repeat_threshold_up:                      viper.GetString("config.buy_repeat_threshold_up"),
		Buy_rsi7_entry:                               viper.GetString("config.buy_rsi7_entry"),
		Buy_MACD_entry:                               viper.GetString("config.buy_MACD_entry"),
		Buy_MACD_upmarket:                            viper.GetString("config.buy_MACD_upmarket"),
		Buy_wait:                                     viper.GetString("config.buy_wait"),
		Exchange_comission:                           viper.GetString("config.exchange_comission"),
		ExchangeName:                                 viper.GetString("config.exchangename"),
		Profit_min:                                   viper.GetString("config.profit_min"),
		SellWaitBeforeCancel:                         viper.GetString("config.sellwaitbeforecancel"),
		SellWaitAfterCancel:                          viper.GetString("config.sellwaitaftercancel"),
		SellToCover:                                  viper.GetString("config.selltocover"),
		Symbol_fiat:                                  viper.GetString("config.symbol_fiat"),
		Symbol_fiat_stash:                            viper.GetString("config.symbol_fiat_stash"),
		Symbol:                                       viper.GetString("config.symbol"),
		Time_enforce:                                 viper.GetString("config.time_enforce"),
		Time_start:                                   viper.GetString("config.time_start"),
		Time_stop:                                    viper.GetString("config.time_stop"),
		TgBotApikey:                                  viper.GetString("config.tgbotapikey"),
		Debug:                                        viper.GetString("config.debug"),
		Exit:                                         viper.GetString("config.exit"),
		DryRun:                                       viper.GetString("config.dryrun"),
		NewSession:                                   viper.GetString("config.newsession"),
		ConfigTemplateList:                           getConfigTemplateList(sessionData),
	}

	return configData

}

/* This routine save viper configuration from html */
func SaveConfigData(
	r *http.Request,
	sessionData *types.Session) {

	viper.Set("config.buy_24hs_highprice_entry", r.PostFormValue("buy_24hs_highprice_entry"))
	viper.Set("config.buy_24hs_highprice_entry_MACD", r.PostFormValue("buy_24hs_highprice_entry_MACD"))
	viper.Set("config.buy_direction_down", r.PostFormValue("buy_direction_down"))
	viper.Set("config.buy_direction_up", r.PostFormValue("buy_direction_up"))
	viper.Set("config.buy_quantity_fiat_up", r.PostFormValue("buy_quantity_fiat_up"))
	viper.Set("config.buy_quantity_fiat_down", r.PostFormValue("buy_quantity_fiat_down"))
	viper.Set("config.buy_quantity_fiat_init", r.PostFormValue("buy_quantity_fiat_init"))
	viper.Set("config.buy_rsi7_entry", r.PostFormValue("buy_rsi7_entry"))
	viper.Set("config.buy_MACD_entry", r.PostFormValue("buy_MACD_entry"))
	viper.Set("config.buy_MACD_upmarket", r.PostFormValue("buy_MACD_upmarket"))
	viper.Set("config.buy_wait", r.PostFormValue("buy_wait"))
	viper.Set("config.buy_repeat_threshold_down", r.PostFormValue("buy_repeat_threshold_down"))
	viper.Set("config.buy_repeat_threshold_down_second", r.PostFormValue("buy_repeat_threshold_down_second"))
	viper.Set("config.buy_repeat_threshold_down_second_start_count", r.PostFormValue("buy_repeat_threshold_down_second_start_count"))
	viper.Set("config.buy_repeat_threshold_up", r.PostFormValue("buy_repeat_threshold_up"))
	viper.Set("config.exchange_comission", r.PostFormValue("exchange_comission"))
	viper.Set("config.exchangename", r.PostFormValue("exchangename"))
	viper.Set("config.profit_min", r.PostFormValue("profit_min"))
	viper.Set("config.sellwaitbeforecancel", r.PostFormValue("sellwaitbeforecancel"))
	viper.Set("config.sellwaitaftercancel", r.PostFormValue("sellwaitaftercancel"))
	viper.Set("config.selltocover", r.PostFormValue("selltocover"))
	viper.Set("config.symbol", r.PostFormValue("symbol"))
	viper.Set("config.symbol_fiat", r.PostFormValue("symbol_fiat"))
	viper.Set("config.symbol_fiat_stash", r.PostFormValue("symbol_fiat_stash"))
	viper.Set("config.time_enforce", r.PostFormValue("time_enforce"))
	viper.Set("config.time_start", r.PostFormValue("time_start"))
	viper.Set("config.time_stop", r.PostFormValue("time_stop"))
	viper.Set("config.debug", r.PostFormValue("debug"))
	viper.Set("config.exit", r.PostFormValue("exit"))
	viper.Set("config.dryrun", r.PostFormValue("dryrun"))
	viper.Set("config.newsession", r.PostFormValue(("newsession")))

	if err := viper.WriteConfig(); err != nil {

		Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			GetFunctionName()+" - "+err.Error())

		os.Exit(1)
	}

}
