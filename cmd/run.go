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

	"github.com/coinpaprika/coinpaprika-api-go-client"
	"github.com/coinpaprika/telegram-bot/telegram"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/telegram-bot-api.v4"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run coinpaprika bot",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}
var debug bool
var token string
var metrics int

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&debug, "debug", "d", false, "enable debugging messages")
	runCmd.Flags().StringVarP(&token, "token", "t", "", "telegram API token")
	runCmd.Flags().IntVarP(&metrics, "metrics", "m", 9900, "metrics port (default :9900) endpoint: /metrics")
	runCmd.MarkFlagRequired("token")
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

			if u.Message == nil {
				log.Debug("Received non-message or non-command")
				continue
			}

			text := `Please use one of the commands:

			/start or /help 	show this message
			/p <symbol> 		check the price for given coin

			/website	 		show link to the coinpaprika webpage
			/source 			show source code of this bot
			`
			log.Debugf("received command: %s", u.Message.Command())
			switch u.Message.Command() {
			case "website":
				text = "https://coinpaprika.com"
			case "source":
				text = "https://github.com/coinpaprika/telegram-bot"
			case "p":
				if text, err = commandP(u.Message.CommandArguments()); err != nil {
					text = err.Error()
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

func commandP(argument string) (string, error) {
	log.Debugf("starting command /p with argument :%s", argument)

	paprikaClient, err := coinpaprika.NewClient()
	if err != nil {
		return "", errors.Wrap(err, "command /p argument:"+argument)
	}

	result, err := paprikaClient.Search(argument, &coinpaprika.SearchOptions{Categories: "currencies"})
	if err != nil {
		return "", errors.Wrap(err, "command /p argument:"+argument)
	}

	log.Debugf("found %d results for command /p with argument :%s", len(result.Currencies), argument)
	if len(result.Currencies) <= 0 {
		return "", errors.Errorf("%s is invalid coin name|ticker|symbol", argument)
	}

	tickerID := result.Currencies[0].ID
	log.Debugf("best match for command /p with argument :%s is: %s", argument, tickerID)

	ticker, err := paprikaClient.GetTickerByID(tickerID)
	if err != nil {
		return "", errors.Wrap(err, "command /p argument:"+argument)
	}

	return fmt.Sprintf("%s price: %f USD, %f BTC", ticker.Name, *ticker.PriceUSD, *ticker.PriceBTC), nil
}
