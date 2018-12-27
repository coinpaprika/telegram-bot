// Copyright © 2018 coinpaprika.com
//
// Licensed under the Apache License, version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"net/http"

	"github.com/coinpaprika/coinpaprika-api-go-client/coinpaprika"
	"github.com/coinpaprika/telegram-bot/telegram"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	debug   bool
	token   string
	metrics int

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run coinpaprika bot",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	commandsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "coinpaprika",
		Subsystem: "telegram_bot",
		Name:      "commands_proccessed",
		Help:      "The total number of processed commands",
	})
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&debug, "debug", "d", false, "enable debugging messages")
	runCmd.Flags().StringVarP(&token, "token", "t", "", "telegram API token")
	runCmd.Flags().IntVarP(&metrics, "metrics", "m", 9900, "metrics port (default :9900) endpoint: /metrics")
	runCmd.MarkFlagRequired("token")

	prometheus.MustRegister(commandsProcessed)
}

func run() error {
	log.SetLevel(log.ErrorLevel)
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("starting telegram-bot")

	bot, err := telegram.NewBot(telegram.BotConfig{
		Token:          token,
		Debug:          debug,
		UpdatesTimeout: 60,
	})

	if err != nil {
		return err
	}

	updates, err := bot.GetUpdatesChannel()
	if err != nil {
		return err
	}
	go func(updates tgbotapi.UpdatesChannel) {
		for u := range updates {
			log.Debugf("Got message: %v", u)

			if u.Message == nil || !u.Message.IsCommand() {
				log.Debug("Received non-message or non-command")
				continue
			}
			commandsProcessed.Inc()

			text := `Please use one of the commands:

			/start or /help 	show this message
			/p <symbol> 		check the coin price
			/s <symbol> 		check the circulating supply
			/v <symbol> 		check the 24h volume

			/source 			show source code of this bot
			`
			log.Debugf("received command: %s", u.Message.Command())
			switch u.Message.Command() {
			case "source":
				text = "https://github.com/coinpaprika/telegram-bot"
			case "p":
				if text, err = commandPrice(u.Message.CommandArguments()); err != nil {
					text = "invalid coin name|ticker|symbol, please try again"
					log.Error(err)
				}
			case "s":
				if text, err = commandSupply(u.Message.CommandArguments()); err != nil {
					text = "invalid coin name|ticker|symbol, please try again"
					log.Error(err)
				}
			case "v":
				if text, err = commandVolume(u.Message.CommandArguments()); err != nil {
					text = "invalid coin name|ticker|symbol, please try again"
					log.Error(err)
				}
			}

			err := bot.SendMessage(telegram.Message{
				ChatID:    int(u.Message.Chat.ID),
				Text:      text,
				MessageID: u.Message.MessageID,
			})

			if err != nil {
				log.Error(err)
			}
		}

	}(updates)

	log.Debugf("launching metrics endpoints :%d/metrics", metrics)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf(":%d", metrics), http.DefaultServeMux)
}

func commandPrice(argument string) (string, error) {
	log.Debugf("processing command /p with argument :%s", argument)

	ticker, err := getTickerByQuery(argument)
	if err != nil {
		return "", errors.Wrap(err, "command /p")
	}

	priceUSD := ticker.Quotes["USD"].Price
	priceBTC := ticker.Quotes["BTC"].Price
	if ticker.Name == nil || ticker.ID == nil || priceUSD == nil || priceBTC == nil {
		return "", errors.Wrap(errors.New("missing data"), "command /p")
	}

	return fmt.Sprintf("%s price: %f USD, %f BTC \n\n http://coinpaprika.com/coin/%s", *ticker.Name, *priceUSD, *priceBTC, *ticker.ID), nil
}

func commandSupply(argument string) (string, error) {
	log.Debugf("processing command /s with argument :%s", argument)

	ticker, err := getTickerByQuery(argument)
	if err != nil {
		return "", errors.Wrap(err, "command /s")
	}

	if ticker.Name == nil || ticker.ID == nil || ticker.CirculatingSupply == nil {
		return "", errors.Wrap(errors.New("missing data"), "command /s")
	}

	return fmt.Sprintf("%s circulating supply: %d \n\n http://coinpaprika.com/coin/%s", *ticker.Name, *ticker.CirculatingSupply, *ticker.ID), nil
}

func commandVolume(argument string) (string, error) {
	log.Debugf("processing command /v with argument :%s", argument)

	ticker, err := getTickerByQuery(argument)
	if err != nil {
		return "", errors.Wrap(err, "command /v")
	}

	volumeUSD := ticker.Quotes["USD"].Volume24h
	if ticker.Name == nil || ticker.ID == nil || volumeUSD == nil {
		return "", errors.Wrap(errors.New("missing data"), "command /v")
	}

	return fmt.Sprintf("%s 24h volume: %.2f USD \n\n http://coinpaprika.com/coin/%s", *ticker.Name, *volumeUSD, *ticker.ID), nil
}

func getTickerByQuery(query string) (*coinpaprika.Ticker, error) {
	paprikaClient := coinpaprika.NewClient(nil)

	searchOpts := &coinpaprika.SearchOptions{Query: query, Categories: "currencies", Modifier: "symbol_search"}
	result, err := paprikaClient.Search.Search(searchOpts)
	if err != nil {
		return nil, errors.Wrap(err, "query:"+query)
	}

	log.Debugf("found %d results for query :%s", len(result.Currencies), query)
	if len(result.Currencies) <= 0 {
		return nil, errors.Errorf("invalid coin name|ticker|symbol")
	}
	if result.Currencies[0].ID == nil {
		return nil, errors.New("missing id for a coin")
	}

	log.Debugf("best match for query :%s is: %s", query, *result.Currencies[0].ID)

	tickerOpts := &coinpaprika.TickersOptions{Quotes: "USD,BTC"}
	ticker, err := paprikaClient.Tickers.GetByID(*result.Currencies[0].ID, tickerOpts)
	if err != nil {
		return nil, errors.Wrap(err, "query:"+query)
	}

	return ticker, nil
}
