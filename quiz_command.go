package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type QuizCommand struct {
	userId int32
}

func (command QuizCommand) Execute() (string, error) {
	entry, err := FindDictionaryEntryToRepeat(command.userId)

	if err != nil {
		// TODO: Test with no rows
		if err == sql.ErrNoRows {
			return "Well done! All repetitions scheduled by now are completed. You're welcome to workout with /quiz later", nil
		} else {
			return "", err
		}
	}

	data := fmt.Sprintf(`{
		"chat_id": %d,
		"text": %q,
		"reply_markup": {
			"keyboard": [
					["remember fast", "took time to remember", "can't remember"]
			],
			"one_time_keyboard": true,
			"resize_keyboard": true
		}
	}`, command.userId, entry.phrase)

	_, err = http.Post(telegramSendMessageEndpoint(), "application/json", strings.NewReader(data))

	if err != nil {
		return "", err
	}

	return "", nil
}

func telegramSendMessageEndpoint() string {
	telegramBotApiUrl, ok := os.LookupEnv("TELEGRAM_BOT_API_URL")

	if !ok {
		panic("TELEGRAM_BOT_API_URL environment variable not set")
	}

	return fmt.Sprintf("%s/sendMessage", telegramBotApiUrl)
}
