## HOW TO USE

Cryptopump opens in your browse it's first instance. 

### METRICS

- On the top left you it shows the thread name and how many instances are in execution.

- Profit: Shows the Total Profit (Total Profit = Sales - Buys); the Net Profit (Net Profit = Total Profit - Order Differences) where Order Difference is the total difference between each order price and the current pair price for all threads. Another way to understand Net Profit is to look at is as the total profit if all orders were to be closed at that moment in time. Net profit is important because CryptoPump will use Profits to buy orders if the crypto pair goes down in price; finally, the average transaction percentage profit across all present and past running threads.

- Thread Profit: Shows the ToNet Profit (Net Thread Profit = Total Thread Profit - Order Thread Differences) where Order Difference is the total difference between each order price and the current pair price for the current threads.; finally, the average transaction percentage profit across the running thread.

- Diff: Shows the sum of Order Differences for the current thread. Order Difference is the total difference between each order price and the current pair price for the current threads.

- Deployed: Shows how much fiat currency is in use across all threads.

- Funds: Shows the total amount of crypto pairs acquired by the current thread and the amount of FIAT currency available for additional purchases.

- Offset : In rare circumstances the database may become out-of-sync with the amount of crypto invested due to the Exchange or Connectivity error. This  field represent the disparity between the system and the exchange quantities. (0 means no difference and all is good)

- Transact./h: Number of Sale transactions per hour.

- MACD is the Moving Average Convergence Divergence it is a trend-following momentum indicator that shows the relationship between two moving averages.

- RSI 14/7/3 is the Relative Strength Index it's an indicator based on closing prices over a duration of specific time.

- Direction: Updated every second from the exchange and is increased at each movement in the same direction, i.e. if the price moves up 10 consecutive times then the direction will be 10.

- Price$: Current price of the selected crypto currency.


## SETTING UP

You can use any template to start your personal configuration and modify the values to suit the trading style you want the bot to have.

### BUY

- Buy Quantity FIAT Upmarket: this is the quantity in fiat currency that the bot will attempt to buy when the direction of the market is up (according to Buy Direction Upmarket), when the pair market price is above the lowest existing transaction waiting to be sold, i.e. 30 for BTC USDT will use $30 USDT every time it buys according to the value of buy direction upmarket. 

- Buy Quantity FIAT Downmarket: this is the quantity in fiat currency that the bot will attempt to buy when the direction of the market is up (according to Buy Direction Downmarket). i.e. 30 for BTC USDT will use $30 every time it buys according to the value of buy direction downmarket.

- Buy Quantity FIAT Initial: The initial amount of fiat currency that should be used by the bot, i.e. if BTC USDT and set to $30 the first buy order will be of $30 USDT.

- Buy Direction Upmarket: This value is the number of consecutive movements the market does increasing the price of the asset before putting a buy order, i.e. if set to 10 the market needs to move up 10 times consecutive before executing a buy order. The higher the value, the more bullish the market needs to be in order to execute a buy order.

- Buy Direction Downmarket: This value is the number of consecutive movement the market does decreasing the price of the asset before putting a buy order, i.e. if set to 10 the market needs to move down 10 times consecutive before executing a buy order.
The higher the value, the more bullish the market needs to be in order to execute a buy order.

- Buy on RSI7: This value indicates the value of RSI (Relative Strength Index) that should be lower so a buy order can be executed, i.e. if set to 45 RSI7 needs to be bellow 45, meaning over sold, so the bot execute the order.

- Buy 24hs HighPrice Entry: This value indicates the maximum amount that bot can buy in relation to the 24 hours highest price, i.e. if set to 0,0003 the bot will buy at a maximum of 0,3% of the highest price.

- Buy Repeat Threshold up: This value indicates the percentage that needs to be increased in relation to the last buy transaction so another buy order is executed, i.e. if set to 0,0001 it needs to be 0,1% different of the previews price.

- Buy Repeat Threshold down: This value indicates the percentage that needs to be decreased in relation to the last buy transaction price so another buy order is executed, i.e. if set to 0,002 it needs to be 2% different of the previews price.

- Buy repeat threshold 2nd down: This value indicates the percentage that needs to be decreased in relation to the first decreased value on a down market so another buy order is executed, i.e. if set to 0,005 it needs to be 5% lower in relation to the previews price. 

- Buy Repeat threshold 2nd Down Start Count: The number os buys in downmarket before "Buy repeat threshold 2nd down" enters into effect.

- Buy Wait: Minimum wait time in seconds before executing buy orders, i.e. if set to 10 it will take 10 seconds between buy orders. 

### SELL

- Minimum Profit: this value indicates the minimum profit so the bot executes a sell order, i.e. if set to 0,005 it will sell an order for 0,5% + exchange commission price. 

- Wait After Cancel: this value indicates the number of seconds the bot waits after canceling an order and performing another, i.e. if set to 10 it will wait 10 seconds to execute another order. 

- Wait Before Cancel: this value indicates the number of seconds the bot waits before canceling an order, i.e. if set to 10 it will wait 10 seconds before cancelling an order. 

- Sell-to-Cover Low Funds: True or False. This option allows the bot to sell your highest buy orders so it can buy at lower values in a down market. (this feature enabled the bot to keep operating, but effectively sell at a loss). 

- Hold Sale on RSI3: This value sets the value of RSI3 (Relative Strength Index) to avoid selling an order, i.e. if set to 70 and RSI3 above 70, meaning the market is over bought, it will not executed a sell order. 

- Stoploss: This option allows the bot to sell your order if the ratio greater than the value, i.e. if the current price is too low compared to the moment it was bought it will sell to avoid increased loss. 

- Exchange Name: the name of the exchange used. Only BINANCE is supported at the moment.

- Exchange commission: The commission taken by the exchange that the bot needs to add when selling an order, i.e. if set to 0,00075 the commission is 0,75% per order when using BNB or set to 0,001 when paying with other currencies for 0,1% commission per order.

- Symbol FIAT: The symbol used for the FIAT currency, i.e. USDT, BUSD. 

- Symbol FIAT Stash: the amount of FIAT that should be preserved by the bot, i.e. if set to 10 and FIAT set to USDT it will always maintain $10 USDT in your trading wallet.

- Symbol: The pair that the bot will trade in this particular instance, i.e. BTCUSDT.

- Enforce Time: True or False, enables the bot to operate during a set period of time set on Start Time and Stop Time. 

- Start Time: If enforce time is set to true this value is used as a start time for the bot operation.

- Stop Time: If enforce time is set to true this value is used to stop the bot operation.

### ORDERS GRID

- OrderID: this value is provided by the exchange when a buy order takes place.

- Quantity: this value indicates the transaction quantity in crypto-currency.

- Quote: this value indicates the transaction total amount in FIAT.

- Price: this value indicates the transaction crypto-currency execution price.

- Target: this value indicates the target crypto-currency price before selling (includes profit margin and exchange fees).

- Diff: this value indicates the difference between current transaction sale price and zero margin sale (includes exchange fee).

- Action [Sell]: The Sell button allows you to execute a sale of an existing order. The 'Diff' lets you know if the order is likely to have a profitable sale or if it is underwater. The sale will occur on the spot market at current market prices.

## STATUS

In the bottom right corner the system status is displayed:

- Buy: Reason in the decision tree on why a given Buy order is not being executed. This field is important and provide information on what configuration tunning might be required.

- Sell: Reason in the decision tree on why a given Sell order is not being executed. This field is important and provide information on what configuration tunning might be required.

- Ops/dec: Number of operation per second. This number is dictated by the crypto-pair volume. Cryptopump analyses every Exchange kline block.

- Signal: Average latency between Cryptopump and the exchange measured every five seconds (best kept below 200ms).

### OTHERS:

- Debug: True or False, enable debug mode output on logs. 

- Exit: True or false, when set to true the bot will stop buying, and when there are no more transactions to be done, meaning all previews buy orders are sold, it will close the instance the bot is running on. 

- DryRun: True or False, when enabled run the bot in DryRun mode without executing a buy or sell order. 

- New Session: True or False, when enabled forces a new session with the bot. Use it if you want to change the Symbol FIAT and Symbol of the trading pair. 

- TestNet: True or False, when enabled starts the bot on Binance TestNet without using real money (require Binance TestNet API keys). 

- Template: select which template the bot will use to avoid writing the same settings multiple times. 


### BUTTONS:

- New: When a session is already in progress it will start a new session on a different HTTP port, i.e. if running the first session on 8080 it will start the next one on 8081. 

- Start: Start the bot on the trading pair previously set. 

- Stop: Stop the bot without selling your active orders. 

- Update: write the changes made within the webui into the configuration file. 

- Buy market: Buy order. The purchase will occur on the spot market at current market prices.

- Sell market: Sell the top order in the orders table. The sale will occur on the spot market at current market prices.


### TELEGRAM:

Telegram allows you to remote monitor that status of your running cryptopump instances, and BUY/SELL orders. The currently available command are:

![](https://github.com/aleibovici/img/blob/b2c9390494906b8e83635a5f320dd48f67a48fbd/telegram_screenshot.jpg?raw=true)

- /report: Provides Available Funds, Deployed Funds, Profit, Return on Investment, Net Profit, Net Return on Investment, Avg. Transaction Percentage gain, Thread Count, System Status, and Master Node.
- /buy: Buy at the current Master Node thread
- /sell: Sell at the current Master Node thread

## RESUMING AND TROUBLESHOOTING:

If you want to stop buy don't want to sell your orders, press stop at each instance. 
To resume start the bot, access the first WebUI, i.e. port 8080, press start. To access the other trading pairs, press new, start the new webui, i.e. port 8081 and press start. Repeat until all instances are resumed. 

If resuming a thread/instance does not work, go into the cryptopump folder and delete the .lock files. Those files are present while the bot is running, if it crashes those won't be deleted so those need to be manually removed before starting the resume process.
