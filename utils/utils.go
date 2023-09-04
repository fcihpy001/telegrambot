package utils

import (
	"encoding/json"
	"log"
	"os"
	"telegramBot/model"
)

func Json2Button(file string, models *[]model.ButtonInfo) {
	path, err := os.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(path, &models)
	if err != nil {
		log.Panic(err)
	}
}
