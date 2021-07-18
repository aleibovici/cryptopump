# cryptopump

![](https://github.com/aleibovici/img/blob/main/cryptopump_screen.png?raw=true)

- CryptoPump is a cryptocurrency trading bot that focuses on extremely high speed and flexibility. The algorithms utilize Go Language and the exchange WebSockets to react in real-time to market movements based on Bollinger statistical analysis and pre-defined profit margins.

- CryptoPump calculates the Relative Strength Index (3,7,14) and MACD index, allowing users to configure buying, selling, and holding thresholds.

- CryptoPump also provides different configuration settings for operating in downmarket, such as specifying the amount to buy in the downmarket when to change purchase behavior and thresholds.

- CryptoPump supports all cryptocurrency pairs and provides the ability to define the exchange commission when calculating profit and when to sell.

- CryptoPump also provides DryRun mode, the ability to use Binance TestNet for testing, Telegram bot integration, Time enforcement, Sell-to-cover, and much more.

- Currently, only the Binance API is supported, but I developed the software to allow easy implementation of additional exchanges.

- Configure the Binance exchange APIKEY and SECRETKEY in config.yml. In addition, the Telegram APIKEY, if in use, has to be configured at TGBOTAPIKEY in the config.yml file.

- Telegram accepts command /stop /sell /buy /funds /master

- CryptoPump requires MySQL to persist data and transactions, and the .sql file to create the structure can be found in the MySQL folder (cryptopump.sql). I use MySQL with Docker in the same machine Cryptopump is running, and it performs well. Cloud-based MySQL instances are also supported. The environment variables are in launch.json if Visual Studio Code is in use; optionally, the following environment variables set DB_USER, DB_PASS, DB_TCP_HOST, DB_PORT, DB_NAME.

- To use Binance TestNet, set launch.json or environment variable TESTNET to True. (https://testnet.binance.vision)

- I run CryptoPump in Visual Studio Code, but it can be run without an IDE. For each instance of the code, a new HTTP port is opened, starting with 8080, 8081, 8082. Just point your browser to the address, and you should get the session configuration page and the Bollinger and Exchange data.

*** If you feel like contributing to the project, you are very welcome ***