# cryptopump

- CryptoPump is a cryptocurrency trading bot that focuses on high speed and flexibility. The algorithms utilize Go Language and exchange WebSockets to react in real-time to market movements based on Bollinger statistical analysis and pre-defined profit margins.

- CryptoPump calculates the Relative Strength Index (3,7,14) and MACD index, allowing users to configure buying, selling, and holding thresholds.

- CryptoPump also provides different configuration settings for operating in the downmarket, such as specifying the amount to buy in the downmarket when to change thresholds.

- CryptoPump supports all cryptocurrency pairs and provides the ability to define the exchange commission when calculating profit and when to sell.

- CryptoPump also provides DryRun mode, the ability to use Binance TestNet for testing, Telegram bot integration, Time enforcement, Sell-to-cover lack of funds, and more.

- Only Binance API is supported, but I developed the software to allow easy implementation of additional exchanges.

- Configure the exchange APIKEY and SECRETKEY in config.yml. Telegram APIKEY, if in use, should be configured at TGBOTAPIKEY.

- Telegram accepts command /stop /sell /buy /funds /master

- CryptoPump requires MySQL to persist data, and the .sql file to create the structure can be found in the MySQL folder (cryptopump.sql). I use MySQL with Docker in the same machine Cryptopump is running, and it works well. Cloud-based MySQL instances are also supported. The environment variables are launch.json if Visual Studio Code is in use; optionally set DB_USER, DB_PASS, DB_TCP_HOST, DB_PORT, DB_NAME.

- To use Binance TestNet set launch.json or environment variable TESTNET to True. (https://testnet.binance.vision)

- I recommend running the code in Visual Studio Code, but it can be run without an IDE. For each instance of the code, a new HTTP port is opened, starting with 8080, 8081, 8082. Just point your browser to the address, and you should get the session configuration and the Bollinger and Exchange data.

*** If you feel like contributing to the project, you are very welcome ***