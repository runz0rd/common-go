package common

import (
	"io/ioutil"
	"time"

	"gopkg.in/tucnak/telebot.v2"
	"gopkg.in/yaml.v3"
)

type TelegramConfig struct {
	BotToken string `yaml:"bot_token,omitempty"`
	UserId   int    `yaml:"user_id,omitempty"`
}

func NewTelegramBot(botToken string) (*telebot.Bot, error) {
	tb, err := telebot.NewBot(telebot.Settings{
		Token:     botToken,
		Poller:    &telebot.LongPoller{Timeout: 10 * time.Second},
		ParseMode: telebot.ModeMarkdown,
	})
	go tb.Start()
	return tb, err
}

func ReadTelegramConfig(path string) (*TelegramConfig, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c TelegramConfig
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
