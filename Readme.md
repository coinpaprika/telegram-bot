# Coinpaprika telegram bot

## Building project

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
