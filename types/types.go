package types

import (
	"database/sql"
	"time"

	"github.com/adshao/go-binance/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sdcoffey/techan"
)

type Order struct {
	ClientOrderID            string  `json:"clientOrderId"`
	CummulativeQuoteQuantity float64 `json:"cummulativeQuoteQty"`
	ExecutedQuantity         float64 `json:"executedQty"`
	OrderID                  int     `json:"orderId"`
	Price                    float64 `json:"price"`
	Side                     string  `json:"side"`
	Status                   string  `json:"status"`
	Symbol                   string  `json:"symbol"`
	TransactTime             int64   `json:"transactTime"`
	ThreadID                 int
	ThreadIDSession          int
}

type Kline struct {
	OpenTime int64  `json:"openTime"`
	Open     string `json:"open"`
	High     string `json:"high"`
	Low      string `json:"low"`
	Close    string `json:"close"`
	Volume   string `json:"volume"`
}

/* WsKline define websocket kline */
type WsKline struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"` /* Currently not in use */
	Symbol               string `json:"s"` /* Currently not in use */
	Interval             string `json:"i"` /* Currently not in use */
	FirstTradeID         int64  `json:"f"` /* Currently not in use */
	LastTradeID          int64  `json:"L"` /* Currently not in use */
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"` /* Currently not in use */
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"` /* Currently not in use */
	ActiveBuyVolume      string `json:"V"` /* Currently not in use */
	ActiveBuyQuoteVolume string `json:"Q"` /* Currently not in use */
}

/* PriceChangeStats define price change stats */
type PriceChangeStats struct {
	HighPrice string `json:"highPrice"`
	LowPrice  string `json:"lowPrice"`
}

type ExchangeInfo struct {
	MaxQuantity string `json:"maxQty"`
	MinQuantity string `json:"minQty"`
	StepSize    string `json:"stepSize"`
}

/* struct for session elements */
type Session struct {
	ThreadID             string
	ThreadIDSession      string
	ThreadCount          int
	SellTransactionCount float64 /* Number of SELL transactions in the last 60 minutes */
	Symbol               string
	Symbol_fiat          string
	Symbol_fiat_funds    float64
	LastBuyTransactTime  time.Time /* This session variable stores the time of the last buy */
	LastSellCanceledTime time.Time /* This session variable stores the time of the cancelled sell */
	ConfigTemplate       int
	ForceBuy             bool             /* This boolean when True force BUY transaction */
	ForceSell            bool             /* This boolean when True force SELL transaction */
	ListenKey            string           /* Listen key for user stream service */
	MasterNode           bool             /* This boolean is true when Master Node is elected */
	TgBotAPI             *tgbotapi.BotAPI /* This variable holds Telegram session bot */
	Db                   *sql.DB          /* mySQL database connection */
	Clients              Client           /* Binance client connection */
	KlineData            []KlineData      /* kline data format for go-echart plotter */
	StopWs               bool             /* Control when to stop Ws Channels */
	Busy                 bool             /* Control wether buy/selling to allow graceful session exit */
	MinQuantity          float64          /* Defines the minimum quantity allowed by exchange */
	MaxQuantity          float64          /* Defines the maximum quantity allowed by exchange */
	StepSize             float64          /* Defines the intervals that a quantity can be increased/decreased by exchange */
}

/* struct for client libraries */
type Client struct {
	Binance *binance.Client
}

/* struct for websocket handlers for exchanges */
type WsHandler struct {
	BinanceWsKline         func(event *binance.WsKlineEvent)      /* WsKlineServe serve websocket kline handler */
	BinanceWsBookTicker    func(event *binance.WsBookTickerEvent) /* WsBookTicker serve websocket kline handler */
	BinanceWsUserDataServe func(message []byte)                   /* WsUserDataServe serve user data handler with listen key */
}

type KlineData struct {
	Date    int64
	Data    [4]float64
	Volumes float64
}

type Market struct {
	Rsi3                      float64            /* Relative Strength Index for 3 periods */
	Rsi7                      float64            /* Relative Strength Index for 7 periods */
	Rsi14                     float64            /* Relative Strength Index for 14 periods */
	MACD                      float64            /* Moving Average for 14 periods */
	Price                     float64            /* Market Price */
	PriceChangeStatsHighPrice float64            /* High price for 1 period */
	PriceChangeStatsLowPrice  float64            /* Low price for 1 period */
	Direction                 int                /* Market Direction */
	TimeStamp                 time.Time          /* Time of last retrieved market Data */
	Series                    *techan.TimeSeries /* kline data format for technical analysis */
}

type Config struct {
	ThreadID                                     interface{} /* For index.html population */
	Apikey                                       interface{} /* Exchange API Key */
	Secretkey                                    interface{} /* Exchange Secret Key */
	ApikeyTestNet                                interface{} /* API key for exchange test network, used with launch.json */
	SecretkeyTestNet                             interface{} /* Secret key for exchange test network, used with launch.json */
	Buy_24hs_highprice_entry                     interface{}
	Buy_24hs_highprice_entry_MACD                interface{}
	Buy_direction_down                           interface{}
	Buy_direction_up                             interface{}
	Buy_quantity_fiat_up                         interface{}
	Buy_quantity_fiat_down                       interface{}
	Buy_quantity_fiat_init                       interface{}
	Buy_repeat_threshold_down                    interface{}
	Buy_repeat_threshold_down_second             interface{}
	Buy_repeat_threshold_down_second_start_count interface{}
	Buy_repeat_threshold_up                      interface{}
	Buy_rsi7_entry                               interface{}
	Buy_MACD_entry                               interface{}
	Buy_wait                                     interface{} /* Wait time between BUY transactions in seconds */
	Exchange_comission                           interface{}
	Profit_min                                   interface{}
	SellWaitBeforeCancel                         interface{} /* Wait time before cancelling a sale in seconds */
	SellWaitAfterCancel                          interface{} /* Wait time before selling after a cancel in seconds */
	SellToCover                                  interface{} /* Define if will sell to cover low funds */
	SellHoldOnRSI3                               interface{} /* Hold sale if RSI3 above defined threshold */
	Symbol_fiat                                  interface{}
	Symbol_fiat_stash                            interface{}
	Symbol                                       interface{}
	Time_enforce                                 interface{}
	Time_start                                   interface{}
	Time_stop                                    interface{}
	Debug                                        interface{}
	Exit                                         interface{}
	DryRun                                       interface{} /* Dry Run mode */
	NewSession                                   interface{} /* Force a new session instead of resume */
	ConfigTemplateList                           interface{} /* List of configuration templates available in ./config folder */
	ExchangeName                                 interface{} /* Exchange name */
	TgBotApikey                                  interface{} /* Telegram bot API key */
	HtmlSnippet                                  interface{} /* Store kline plotter graph for html output */
	Orders                                       interface{} /* Store thread orders for html output */
	FiatFunds                                    interface{} /* Store fiat currency funds for html output */
	Profit                                       interface{} /* Store total profit for html output */
	ProfitThreadID                               interface{} /* Store threadID profit for html output */
	SellTransactionCount                         interface{} /* Store Number of SELL transactions in the last 60 minutes for html output */
	ThreadCount                                  interface{} /* Store thread count for html output */
	ThreadAmount                                 interface{} /* Store thread cost amount for html output */
	MarketDataMACD                               interface{} /* Store MACD for html output */
	MarketDataRsi3                               interface{} /* Store RSI14 for html output */
	MarketDataRsi7                               interface{} /* Store RSI7 for html output */
	MarketDataRsi14                              interface{} /* Store RSI3 for html output */
}

/* Struct for User Data Streams for Binance */
type OutboundAccountPosition struct {
	EventType  string     `json:"e"` /* Event type */
	EventTime  int64      `json:"E"` /* Event Time */
	LastUpdate int64      `json:"u"` /* Time of last account update */
	Balances   []Balances `json:"B"` /* Balances Array */
}

/* Struct for User Data Streams for Binance */
type Balances struct {
	Asset  string `json:"a"` /* Asset */
	Free   string `json:"f"` /* Free */
	Locked string `json:"l"` /* Locked */
}

type ExecutionReport struct {
	EventType             string `json:"e"` //Event type
	EventTime             int64  `json:"E"` //Event Time
	Symbol                string `json:"s"` //Symbol
	ClientOrderID         string `json:"c"` //Client order ID
	Side                  string `json:"S"` //Side
	OrderType             string `json:"o"` //Order type
	TimeInForce           string `json:"f"` //Time in force
	Quantity              string `json:"q"` //Order quantity
	Price                 string `json:"p"` //Order price
	StopPrice             string `json:"P"` //Stop price
	IcebergQuantity       string `json:"F"` //Iceberg quantity
	OrderListId           int64  `json:"g"` //OrderListId
	OriginalClientOrderID string `json:"C"` //Original client order ID; This is the ID of the order being canceled
	ExecutionType         string `json:"x"` //Current execution type
	Status                string `json:"X"` //Current order status
	OrderRejectReason     string `json:"r"` //Order reject reason; will be an error code.
	OrderID               int    `json:"i"` //Order ID
	LastExecutedQuantity  string `json:"l"` //Last executed quantity
	CummulativeQty        string `json:"z"` //Cumulative filled quantity
	LastExecutedPrice     string `json:"L"` //Last executed price
	ComissionAmount       string `json:"n"` //Commission amount
	ComissionAsset        string `json:"N"` //Commission asset
	TransactTime          int64  `json:"T"` //Transaction time
	TradeID               int    `json:"t"` //Trade ID
	Ignore0               int    `json:"I"` //Ignore
	IsOrderOnTheBook      bool   `json:"w"` //Is the order on the book?
	IsTradeMakerSide      bool   `json:"m"` //Is this trade the maker side?
	Ignore1               bool   `json:"M"` //Ignore
	OrderCreationTime     int64  `json:"O"` //Order creation time
	CummulativeQuoteQty   string `json:"Z"` //Cumulative quote asset transacted quantity
	LastQuoteQty          string `json:"Y"` //Last quote asset transacted quantity (i.e. lastPrice * lastQty)
	QuoteOrderQty         string `json:"Q"` //Quote Order Qty
}
