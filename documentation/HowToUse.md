## HOW TO USE

Cryptopump opens in your browse it's first instance. 

On the top left you it shows the thread name and how many instances are running. 
Profit shows the profit in the currency set and it's %.*
Deployed show how much fiat currency is used.*
Funds shows the total amount of fiat you have available in the selected trading pair. 
Transact./h shows how many transaction the bot did per hour. 

*If more than one instance is running the total amount will be updated.

On the top right you can see indicators for that particular trading pair.
MACD is the Moving Average Convergence Divergence it is a trend-following momentum indicator that shows the relationship between two moving averages.
RSI 14/7/3 is the Relative Strengfht Index it's an indicator based on closing prices over a duration of specific time.
Direction is updated every second from your exchange and is greater at each movement in the same direction, i.e. if the price moves up 10 consecutive times it will show as 10 under direction.
Price$ is the cost of the selected currency.

## SETTING UP

You can use any template to start your personal configuration and modify the values to suit the trading style you want the bot to have.

### BUY

- Buy Quantity FIAT Upmarket: this is the quantity in fiat currency that the bot will attempt to buy when the direction of the market is up (according to Buy Direction Upmarket), when the pair market price is above the lowest existing transaction waiting to be sold, i.e. 30 for BTC USDT will use $30 USDT every time it buys according to the value of buy direction upmarket. 

- Buy Quantity FIAT Downmarket: this is the quantity in fiat currency that the bot will attempt to buy when the direction of the market is up (according to Buy Direction Downmarket). i.e. 30 for BTC USDT will use $30 every time it buys according to the value of buy direction downmarket.

- Buy Quantity FIAT Initial: The initial amount of fiat currency that should be used by the bot, i.e. if BTC USDT and set to $30 the first buy order will be of $30 USDT.

- Buy Direction Upmarket: This value is the number of consecutive movements the market does increasing the price of the asset before putting a buy order, i.e. if set to 10 the market needs to move up 10 times consecutive before executing a buy order. The higher the value, the more bulish the market needs to be in order to execute a buy order.

- Buy Direction Downmarket: This value is the number of consecutive movement the market does decreasing the price of the asset before putting a buy order, i.e. if set to 10 the market needs to move down 10 times consective before executing a buy order.
The higher the value, the more bulish the market needs to be in order to execute a buy order.

- Buy on RSI7: This value indicates the value of RSI (Relative Strengh Index) that should be lower so a buy order can be executed, i.e. if set to 45 RSI7 needs to be bellow 45, meaning over sold, so the bot execute the order.

- Buy 24hs HighPrice Entry: This value indicates the maximum amount that bot can buy in relation to the 24 hours highest price, i.e. if set to 0,0003 the bot will buy at a maximum of 0,3% of the highest price.

- Buy Repeat Threshold up: This value indicates the percentage that needs to be increased in relation to the last buy transaction so another buy order is executed, i.e. if set to 0,0001 it needs to be 0,1% different of the previews price.

- Buy Repeat Thresold down: This value indicates the percentage that needs to be decreased in relation to the last buy transaction price so another buy order is executed, i.e. if set to 0,002 it needs to be 2% different of the previews price.

- Buy repeat threshold 2nd down: This value indicates the percentage that needs to be decreased in relation to the first decreased value on a down market so another buy order is executed, i.e. if set to 0,005 it needs to be 5% lower in relation to the previews price. 

- Buy Repeat threashold 2nd Down Start Count: The number os buys in downmarket before "Buy repeat threshold 2nd down" enters into effect.

- Buy Wait: Minimum wait time in seconds before executing buy orders, i.e. if set to 10 it will take 10 seconds between buy orders. 

### SELL

- Minimum Profit: this value indicates the minimum profit so the bot executes a sell order, i.e. if set to 0,005 it will sell an order for 0,5% + exchange comission price. 

- Wait After Cancel: this value indicates the number of seconds the bot waits after canceling an order and performing another, i.e. if set to 10 it will wait 10 seconds to execute another order. 

- Wait Before Cancel: this value indicates the number of seconds the bot waits before canceling an order, i.e. if set to 10 it will wait 10 seconds before cancelling an order. 

- Sell-to-Cover Low Funds: True or False. This option allows the bot to sell your highest buy orders so it can buy at lower values in a down market. (this feature enabled the bot to keep operating, but effectivile sell at a loss). 

- Hold Sale on RSI3: This value sets the value of RSI3 (Relative Strenght Index) to avoid selling an order, i.e. if set to 70 and RSI3 above 70, meaning the market is over bought, it will not executed a sell order. 

- Stoploss: This option allows the bot to sell your order if the ratio greater than the value, i.e. if the current price is too low compared to the moment it was bought it will sell to avoid increased loss. 

- Exchange Name: the name of the exchange used. Only BINANCE is supported at the moment.

- Exchange commission: The comission taken by the exchange that the bot needs to add when selling an order, i.e. if set to 0,00075 the comission is 0,75% per order when using BNB or set to 0,001 when paying with other currencies for 0,1% comission per order.

- Symbol FIAT: The symbol used for the FIAT currency, i.e. USDT, BUSD. 

- Symbol FIAT Stash: the amount of FIAT that should be preserved by the bot, i.e. if set to 10 and FIAT set to USDT it will always maintein $10 USDT in your trading wallet.

- Symbol: The pair that the bot will trade in this particular instance, i.e. BTCUSDT.

- Enforce Time: True or False, enables the bot to operate during a set period of time set on Start Time and Stop Time. 

- Start Time: If enforce time is set to true this value is used as a start time for the bot operation.

- Stop Time: If enforce time is set to true this value is used to stop the bot operation.

### OTHERS:

- Debug: True or False, enable debug mode output on logs. 

- Exit: True or false, when set to true the bot will stop buying, and when there are no more transactions to be done, meaning all previews buy orders are sold, it will close the instance the bot is running on. 

- DryRun: True or False, when enabled run the bot in DryRun mode without executing a buy or sell order. 

- New Session: True or False, when enabled forces a new session with the bot. Use it if you want to change the Symbol FIAT and Symbol of the trading pair. 

- TestNet: True or False, when enabled starts the bot on Binance TestNet without using real money (reques Binance TestNet API keys). 

- Template: select which template the bot will use to avoid writting the same settings multiple times. 


### BUTTONS:

- New: When a session is already in progress it will start a new session on a different HTTP port, i.e. if running the first session on 8080 it will start the next one on 8081. 

- Start: Start the bot on the trading pair previewsly set. 

- Stop: Stop the bot without selling your active orders. 

- Update: write the changes made within the webui into the configuration file. 

- Buy market: executes a buy order at the price in that particular moment. 

- Sell market: executes a sell order at the price in that particular moment, each press will sell one particular order, press multiple times to sell all. 



## RESUMING AND TROUBLESHOOTING:

If you want to stop buy don't want to sell your orders, press stop at each instance. 
To resume start the bot, access the first webui, i.e. port 8080, press start. To access the other trading pairs, press new, start the new webui, i.e. port 8081 and press start. Repeat until all instances are resumed. 

If resuming a thread/instance does not work, go into the cryptopump folder and delete the .lock files. Those files are present while the bot is running, if it crashes those won't be deleted so those need to be manually removed before starting the resume process. 