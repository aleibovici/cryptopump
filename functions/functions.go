package functions

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"os"
	"strconv"

	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/types"
	"github.com/tcnksm/go-httpstat"

	"github.com/rs/xid"
)

// StrToFloat64 function
/* This public function convert string to float64 */
func StrToFloat64(value string) (r float64) {

	var err error

	if r, err = strconv.ParseFloat(value, 8); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		return 0

	}

	return r
}

// Float64ToStr function
/* This public function convert float64 to string with variable precision */
func Float64ToStr(value float64, prec int) string {

	return strconv.FormatFloat(value, 'f', prec, 64)

}

// IntToFloat64 convert Int to Float64
func IntToFloat64(value int) float64 {

	return float64(value)

}

// StrToInt convert string to int
func StrToInt(value string) (r int) {

	var err error

	if r, err = strconv.Atoi(value); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	return r

}

// MustGetenv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func MustGetenv(k string) string {

	v := os.Getenv(k)
	if v == "" {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName() + " - " + "Environment variable not set",
			LogLevel: "DebugLevel",
		}.Do()

	}

	return strings.ToLower(v)

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

// GetThreadID Return random thread ID
func GetThreadID() string {

	return xid.New().String()

}

/* Convert Strign to Time */
func stringToTime(str string) (r time.Time) {

	var err error

	if r, err = time.Parse(time.Kitchen, str); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	return r

}

// IsInTimeRange Check if time is in a specific range
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

// IsFundsAvailable Validate available funds to buy
func IsFundsAvailable(
	configData *types.Config,
	sessionData *types.Session) bool {

	return (sessionData.SymbolFiatFunds - configData.SymbolFiatStash) >= configData.BuyQuantityFiatDown

}

/* Select the correct html template based on sessionData */
func selectTemplate(
	sessionData *types.Session) (template string) {

	if sessionData.Admin {

		template = "admin.html" /* Admin template */

	} else if sessionData.ThreadID == "" {

		template = "index.html"

	} else {

		template = "index_nostart.html"

	}

	return template

}

// ExecuteTemplate is responsible for executing any templates
func ExecuteTemplate(
	wr io.Writer,
	data interface{},
	sessionData *types.Session) {

	var tlp *template.Template
	var err error

	if tlp, err = template.ParseGlob("./templates/*"); err != nil {

		defer os.Exit(1)

	}

	if err = tlp.ExecuteTemplate(wr, selectTemplate(sessionData), data); err != nil {

		defer os.Exit(1)

	}

	/* Conditional defer logging when there is an error retriving data */
	defer func() {
		if err != nil {
			logger.LogEntry{ /* Log Entry */
				Config:   nil,
				Market:   nil,
				Session:  nil,
				Order:    &types.Order{},
				Message:  GetFunctionName() + " - " + err.Error(),
				LogLevel: "DebugLevel",
			}.Do()
		}
	}()

}

// GetFunctionName Retrieve current function name
func GetFunctionName() string {

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return frame.Function

}

// GetPort Determine port for HTTP service.
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

// DeleteConfigFile Delete configuration file for ThreadID
func DeleteConfigFile(sessionData *types.Session) {

	filename := sessionData.ThreadID + ".yml"
	path := "./config/"

	if err := os.Remove(path + filename); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		return

	}

}

// SaveConfigGlobalData save viper configuration from html
func SaveConfigGlobalData(
	viperData *types.ViperData,
	r *http.Request,
	sessionData *types.Session) {

	viperData.V2.Set("config_global.apiKey", r.FormValue("Apikey"))                     /* Api Key */
	viperData.V2.Set("config_global.secretKey", r.FormValue("Secretkey"))               /* Secret Key */
	viperData.V2.Set("config_global.apiKeyTestNet", r.FormValue("ApikeyTestNet"))       /* Api Key TestNet */
	viperData.V2.Set("config_global.secretKeyTestNet", r.FormValue("SecretkeyTestNet")) /* Secret Key TestNet */
	viperData.V2.Set("config_global.tgbotapikey", r.FormValue("TgBotApikey"))           /* Tg Bot Api Key */

	if err := viperData.V2.WriteConfig(); err != nil { /* Write configuration file */

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	logger.LogEntry{ /* Log Entry */
		Config:   nil,
		Market:   nil,
		Session:  sessionData,
		Order:    &types.Order{},
		Message:  "Global configuration saved",
		LogLevel: "InfoLevel",
	}.Do()

}

// GetConfigData Retrieve or create config file based on ThreadID
func GetConfigData(
	viperData *types.ViperData,
	sessionData *types.Session) *types.Config {

	configData := loadConfigData(viperData, sessionData)

	if sessionData.ThreadID != "" {

		filename := sessionData.ThreadID + ".yml"
		writePath := "./config/"

		if _, err := os.Stat(writePath + filename); err == nil {

			/* Test for existing ThreadID config file and load configuration */
			viperData.V1.SetConfigFile(writePath + filename)

			if err := viperData.V1.ReadInConfig(); err != nil {

				logger.LogEntry{ /* Log Entry */
					Config:   nil,
					Market:   nil,
					Session:  nil,
					Order:    &types.Order{},
					Message:  GetFunctionName() + " - " + err.Error(),
					LogLevel: "DebugLevel",
				}.Do()

			}

			configData = loadConfigData(viperData, sessionData)

		} else if os.IsNotExist(err) {

			/* Create new ThreadID config file and load configuration */
			if err := viperData.V1.WriteConfigAs(writePath + filename); err != nil {

				logger.LogEntry{ /* Log Entry */
					Config:   nil,
					Market:   nil,
					Session:  nil,
					Order:    &types.Order{},
					Message:  GetFunctionName() + " - " + err.Error(),
					LogLevel: "DebugLevel",
				}.Do()

			}

			viperData.V1.SetConfigFile(writePath + filename)

			if err := viperData.V1.ReadInConfig(); err != nil {

				logger.LogEntry{ /* Log Entry */
					Config:   nil,
					Market:   nil,
					Session:  nil,
					Order:    &types.Order{},
					Message:  GetFunctionName() + " - " + err.Error(),
					LogLevel: "DebugLevel",
				}.Do()

			}

			configData = loadConfigData(viperData, sessionData)

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

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

		os.Exit(1)

	}

	return files

}

// LoadConfigTemplate Load the selected configuration template
// Three's a BUG where it only works before the first UPDATE
func LoadConfigTemplate(
	viperData *types.ViperData,
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

	filenameOld := viperData.V1.ConfigFileUsed()

	/* Set selected template as current config and load settings and return configData*/
	viperData.V1.SetConfigFile("./config/" + filename)
	if err := viperData.V1.ReadInConfig(); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()
	}

	configData := loadConfigData(viperData, sessionData)

	/* Set origina template as current config */
	viperData.V1.SetConfigFile(filenameOld)
	if err := viperData.V1.ReadInConfig(); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	return configData

}

/* This routine load viper configuration data into map[string]interface{} */
func loadConfigData(
	viperData *types.ViperData,
	sessionData *types.Session) *types.Config {

	configData := &types.Config{
		ThreadID:                               sessionData.ThreadID,
		Buy24hsHighpriceEntry:                  viperData.V1.GetFloat64("config.buy_24hs_highprice_entry"),
		BuyDirectionDown:                       viperData.V1.GetInt("config.buy_direction_down"),
		BuyDirectionUp:                         viperData.V1.GetInt("config.buy_direction_up"),
		BuyQuantityFiatUp:                      viperData.V1.GetFloat64("config.buy_quantity_fiat_up"),
		BuyQuantityFiatDown:                    viperData.V1.GetFloat64("config.buy_quantity_fiat_down"),
		BuyQuantityFiatInit:                    viperData.V1.GetFloat64("config.buy_quantity_fiat_init"),
		BuyRepeatThresholdDown:                 viperData.V1.GetFloat64("config.buy_repeat_threshold_down"),
		BuyRepeatThresholdDownSecond:           viperData.V1.GetFloat64("config.buy_repeat_threshold_down_second"),
		BuyRepeatThresholdDownSecondStartCount: viperData.V1.GetInt("config.buy_repeat_threshold_down_second_start_count"),
		BuyRepeatThresholdUp:                   viperData.V1.GetFloat64("config.buy_repeat_threshold_up"),
		BuyRsi7Entry:                           viperData.V1.GetFloat64("config.buy_rsi7_entry"),
		BuyWait:                                viperData.V1.GetInt("config.buy_wait"),
		ExchangeComission:                      viperData.V1.GetFloat64("config.exchange_comission"),
		ProfitMin:                              viperData.V1.GetFloat64("config.profit_min"),
		SellWaitBeforeCancel:                   viperData.V1.GetInt("config.sellwaitbeforecancel"),
		SellWaitAfterCancel:                    viperData.V1.GetInt("config.sellwaitaftercancel"),
		SellToCover:                            viperData.V1.GetBool("config.selltocover"),
		SellHoldOnRSI3:                         viperData.V1.GetFloat64("config.sellholdonrsi3"),
		Stoploss:                               viperData.V1.GetFloat64("config.stoploss"),
		SymbolFiat:                             viperData.V1.GetString("config.symbol_fiat"),
		SymbolFiatStash:                        viperData.V1.GetFloat64("config.symbol_fiat_stash"),
		Symbol:                                 viperData.V1.GetString("config.symbol"),
		TimeEnforce:                            viperData.V1.GetBool("config.time_enforce"),
		TimeStart:                              viperData.V1.GetString("config.time_start"),
		TimeStop:                               viperData.V1.GetString("config.time_stop"),
		Debug:                                  viperData.V1.GetBool("config.debug"),
		Exit:                                   viperData.V1.GetBool("config.exit"),
		DryRun:                                 viperData.V1.GetBool("config.dryrun"),
		NewSession:                             viperData.V1.GetBool("config.newsession"),
		ConfigTemplateList:                     getConfigTemplateList(sessionData),
		ExchangeName:                           viperData.V1.GetString("config.exchangename"),
		TestNet:                                viperData.V1.GetBool("config.testnet"),
		HTMLSnippet:                            nil,
		ConfigGlobal: &types.ConfigGlobal{
			Apikey:           viperData.V2.GetString("config_global.apiKey"),
			Secretkey:        viperData.V2.GetString("config_global.secretKey"),
			ApikeyTestNet:    viperData.V2.GetString("config_global.apiKeyTestNet"),
			SecretkeyTestNet: viperData.V2.GetString("config_global.secretKeyTestNet"),
			TgBotApikey:      viperData.V2.GetString("config_global.tgbotapikey")},
	}

	return configData

}

// SaveConfigData save viper configuration from html
func SaveConfigData(
	viperData *types.ViperData,
	r *http.Request,
	sessionData *types.Session) {

	viperData.V1.Set("config.buy_24hs_highprice_entry", r.PostFormValue("buy24hsHighpriceEntry"))
	viperData.V1.Set("config.buy_direction_down", r.PostFormValue("buyDirectionDown"))
	viperData.V1.Set("config.buy_direction_up", r.PostFormValue("buyDirectionUp"))
	viperData.V1.Set("config.buy_quantity_fiat_up", r.PostFormValue("buyQuantityFiatUp"))
	viperData.V1.Set("config.buy_quantity_fiat_down", r.PostFormValue("buyQuantityFiatDown"))
	viperData.V1.Set("config.buy_quantity_fiat_init", r.PostFormValue("buyQuantityFiatInit"))
	viperData.V1.Set("config.buy_rsi7_entry", r.PostFormValue("buyRsi7Entry"))
	viperData.V1.Set("config.buy_wait", r.PostFormValue("buyWait"))
	viperData.V1.Set("config.buy_repeat_threshold_down", r.PostFormValue("buyRepeatThresholdDown"))
	viperData.V1.Set("config.buy_repeat_threshold_down_second", r.PostFormValue("buyRepeatThresholdDownSecond"))
	viperData.V1.Set("config.buy_repeat_threshold_down_second_start_count", r.PostFormValue("buyRepeatThresholdDownSecondStartCount"))
	viperData.V1.Set("config.buy_repeat_threshold_up", r.PostFormValue("buyRepeatThresholdUp"))
	viperData.V1.Set("config.exchange_comission", r.PostFormValue("exchangeComission"))
	if r.PostFormValue("exchangename") != "" { /* Test for disabled input in index_nostart.html where return is nil */
		viperData.V1.Set("config.exchangename", r.PostFormValue("exchangename"))
	}
	viperData.V1.Set("config.profit_min", r.PostFormValue("profitMin"))
	viperData.V1.Set("config.sellwaitbeforecancel", r.PostFormValue("sellwaitbeforecancel"))
	viperData.V1.Set("config.sellwaitaftercancel", r.PostFormValue("sellwaitaftercancel"))
	viperData.V1.Set("config.selltocover", r.PostFormValue("selltocover"))
	viperData.V1.Set("config.sellholdonrsi3", r.PostFormValue("sellholdonrsi3"))
	viperData.V1.Set("config.Stoploss", r.PostFormValue("stoploss"))
	if r.PostFormValue("exchangename") != "" { /* Test for disabled input in index_nostart.html where return is nil */
		viperData.V1.Set("config.symbol", r.PostFormValue("symbol"))
	}
	if r.PostFormValue("exchangename") != "" { /* Test for disabled input in index_nostart.html where return is nil */
		viperData.V1.Set("config.symbol_fiat", r.PostFormValue("symbol_fiat"))
	}
	viperData.V1.Set("config.symbol_fiat_stash", r.PostFormValue("symbolFiatStash"))
	viperData.V1.Set("config.time_enforce", r.PostFormValue("timeEnforce"))
	viperData.V1.Set("config.time_start", r.PostFormValue("timeStart"))
	viperData.V1.Set("config.time_stop", r.PostFormValue("timeStop"))
	if r.PostFormValue("exchangename") != "" { /* Test for disabled input in index_nostart.html where return is nil */
		viperData.V1.Set("config.testnet", r.PostFormValue("testnet"))
	}
	viperData.V1.Set("config.debug", r.PostFormValue("debug"))
	viperData.V1.Set("config.exit", r.PostFormValue("exit"))
	viperData.V1.Set("config.dryrun", r.PostFormValue("dryrun"))
	if r.PostFormValue("exchangename") != "" { /* Test for disabled input in index_nostart.html where return is nil */
		viperData.V1.Set("config.newsession", r.PostFormValue(("newsession")))
	}

	if err := viperData.V1.WriteConfig(); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  nil,
			Order:    &types.Order{},
			Message:  GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	logger.LogEntry{ /* Log Entry */
		Config:   nil,
		Market:   nil,
		Session:  sessionData,
		Order:    &types.Order{},
		Message:  "Configuration saved",
		LogLevel: "InfoLevel",
	}.Do()

}

// GetExchangeLatency retrieve the latency between the exchange and client
func GetExchangeLatency(sessionData *types.Session) (latency int64, err error) {

	/* Package httpstat traces HTTP latency infomation
	(DNSLookup, TCP Connection and so on) on any golang HTTP request. */

	var req *http.Request
	var res *http.Response

	if req, err = http.NewRequest("GET", sessionData.Clients.Binance.BaseURL, nil); err != nil {

		return 0, err

	}

	/* Create go-httpstat powered context and pass it to http.Request */
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient
	if res, err = client.Do(req); err != nil {

		return 0, err

	}

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	result.End(time.Now())

	return result.ServerProcessing.Milliseconds(), err

}
