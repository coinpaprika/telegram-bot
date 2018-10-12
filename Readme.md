[![Build Status](https://travis-ci.org/coinpaprika/telegram-bot.svg?branch=master)](https://travis-ci.org/coinpaprika/telegram-bot)
[![go-doc](https://godoc.org/github.com/coinpaprika/telegram-bot?status.svg)](https://godoc.org/github.com/coinpaprika/telegram-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/coinpaprika/telegram-bot)](https://goreportcard.com/report/github.com/coinpaprika/telegram-bot)

# Coinpaprika telegram bot

## Commands

```
   /start or /help  show this message
   /p <symbol>      check the coin price
   /s <symbol>      check the circulating supply
   /v <symbol>      check the 24h volume
   /source          show source code of this bot
```   

## Telegram address 
https://t.me/CoinpaprikaBot

## Binary releases
https://github.com/coinpaprika/telegram-bot/releases

## Building project from source

```
git clone git@github.com:coinpaprika/telegram-bot.git
cd telegram-bot/
make 
```

## Running bot
Basic usage: ```./telegram-bot run -t "telegram_bot_api_key"```

Where telegram_bot_api_key can be generated as described https://core.telegram.org/bots#creating-a-new-bot 


Additional parameters are described in help section:
```./telegram-bot run --help```

By default [/metrics](http://localhost:9900/metrics) endopoint is launched which is compatibile with https://prometheus.io/

## Version checking
```./telegram-bot version```
