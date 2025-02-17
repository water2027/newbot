package config

import (
	"os"
	"encoding/json"
	"io"
	"log"
)

type BotConfig struct {
	TargetGroupName []string `json:"targetGroupName"` //目标群名
	Str             string `json:"str"` //格式化输出的字符串
	Url				string `json:"url"` //sse的url
	TimeInterval 	int `json:"timeInterval"`
	Telephone       string `json:"telephone"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	StartNum        int `json:"startNum"`
}

var botConfig BotConfig

func InitBotConfig() {
	jsonFile, err := os.Open("botconfig.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &botConfig)
}

func GetBotConfig() BotConfig {
	return botConfig
}

func UpdateBotConfig(newBotConfig BotConfig) {
	botConfig = newBotConfig
}