# CryptoPump

CryptoPump is a cryptocurrency trading tool that focuses on extremely high speed and flexibility.

[![Go](https://github.com/aleibovici/cryptopump/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/aleibovici/cryptopump/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aleibovici/cryptopump)](https://goreportcard.com/report/github.com/aleibovici/cryptopump)
[![Coverage Status](https://coveralls.io/repos/github/aleibovici/cryptopump/badge.svg?branch=main)](https://coveralls.io/github/aleibovici/cryptopump?branch=main)
[![Codacy Security Scan](https://github.com/aleibovici/cryptopump/actions/workflows/codacy-analysis.yml/badge.svg?branch=main)](https://github.com/aleibovici/cryptopump/actions/workflows/codacy-analysis.yml)
[![Maintainability](https://api.codeclimate.com/v1/badges/62f86b5a3d94b1e2e355/maintainability)](https://codeclimate.com/github/aleibovici/cryptopump/maintainability)

![](https://github.com/aleibovici/img/blob/main/cryptopump_screen.png?raw=true)

Do not risk money which you are afraid to lose. USE THE SOFTWARE AT YOUR OWN RISK. THE AUTHORS AND ALL AFFILIATES ASSUME NO RESPONSIBILITY FOR YOUR TRADING RESULTS.

Always start by running a this trading tool in Dry-run or TestNet and do not engage money before you understand how it works and what profit/loss you should expect.
#### - CryptoPump is now available as a self-contained Docker container set for linux/amd64 and linux/arm/v7 (Raspberry Pi). Check it out at https://hub.docker.com/repository/docker/andreleibovici/cryptopump

- CryptoPump is a cryptocurrency trading tool that focuses on extremely high speed and flexibility. The algorithms utilize Go Language and the exchange WebSockets to react in real-time to market movements based on Bollinger statistical analysis and pre-defined profit margins.

- CryptoPump is easy to deploy and with Docker it can be up and running in minutes.

- CryptoPump calculates the Relative Strength Index (3,7,14), MACD index, and Market Volume Direction, allowing you to configure buying, selling, and holding thresholds.

- CryptoPump also provides different configuration settings for operating in downmarket, such as specifying the amount to buy in the downmarket when to change purchase behavior and thresholds.

- CryptoPump supports all cryptocurrency pairs and provides the ability to define the exchange commission when calculating profit and when to sell.

- CryptoPump also provides DryRun mode, the ability to use Binance TestNet for testing, Telegram bot integration, Time enforcement, Sell-to-cover, and more. (<https://testnet.binance.vision>)

- CryptoPump currently only support Binance API but it was developed to allow easy implementation of additional exchanges.

- CryptoPump has a native Telegram bot that accepts commands /stop /sell /buy and /report. Telegram will also alert you if any issues happen.

![](https://github.com/aleibovici/img/blob/b2c9390494906b8e83635a5f320dd48f67a48fbd/telegram_screenshot.jpg?raw=true)

- CryptoPump requires MySQL to persist data and transactions, and the .sql file to create the structure can be found in the MySQL folder (cryptopump.sql). I use MySQL with Docker in the same machine Cryptopump is running, and it performs well. Cloud-based MySQL instances are also supported. The environment variables are in launch.json if Visual Studio Code is in use; optionally, the following environment variables set DB_USER, DB_PASS, DB_TCP_HOST, DB_PORT, DB_NAME. For using MySQL with docker go here (<https://hub.docker.com/_/mysql>). (refer to HOW TO INSTALL file)

- For each instance of the code, a new HTTP port is opened, starting with 8080, 8081, 8082 (or starting with the port defined by environment variable PORT). Just point your browser to the address, and you should get the session configuration page and the Bollinger and Exchange data.
