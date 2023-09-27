# TonListener
Microservice in Go that tracks the withdraws and deposits of the users and notifies them in Telegram

## Clone the project

```
$ git clone https://github.com/fantarics/ton_listener_service
```

## Environment

```
XTOKEN - bearer token of the tonapi.io service
TON_URL - url for tonapi.io
MAIN_ADDRESS - address of the main hot wallet
TELEGRAM_TOKEN - token of the telegram bot
CONFIG_TON - url to interact with wallets
NOTIFICATION_DESTINATION - notification of the every deposit to telegram

# Database configuration example
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=cap_wallet_db
```

## Connect
You can get list of public lite servers from official TON configs

CONFIG_TON:
* Mainnet - `https://ton.org/global.config.json`
* Testnet - `https://ton-blockchain.github.io/testnet-global.config.json`

TON_URL:
* Mainnet - `https://tonapi.io`
* Testnet - `https://testnet.tonapi.io`

## Setup 

1. Setup the env variables
2. download docker.io & docker-compose
3. Run the command `docker-compose up -d`