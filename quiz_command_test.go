package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecuteQuizCommand(t *testing.T) {
	clearDB()
	entry := DictionaryEntry{userId: 7, phrase: "fusion", translation: "слияние, синтез", phraseMemorizationStep: 0,
		translationMemorizationStep: 0, nextPhraseRepetition: time.Now(), nextTranslationRepetition: time.Now()}
	entry.Save()

	requestHasBeenMade := false
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Method, http.MethodPost, "Expected request method to be POST")
		require.Equal(t, r.Header.Get("Content-Type"), "application/json")

		body_bytes, _ := io.ReadAll(r.Body)
		body := string(body_bytes)
		require.Contains(t, body, fmt.Sprintf(`"chat_id": %d`, entry.userId))
		require.Contains(t, body, fmt.Sprintf(`"text": %q`, entry.phrase))

		w.WriteHeader(http.StatusOK)
		requestHasBeenMade = true
	}))

	defer svr.Close()

	os.Setenv("TELEGRAM_BOT_API_URL", svr.URL)

	command := QuizCommand{userId: 7}
	result, err := command.Execute()

	require.True(t, requestHasBeenMade)
	require.Empty(t, result)
	require.Nil(t, err)

	// Ask to remember translation for the phrase
	// -- Select entries: due_to_remember = min(next_repetition | next_repetition < now)
	// -- If count(due_to_remember) < 1 then Send message "No phrases due to remember"
	// -- Send Telegram message asking to remember selected phrase
	// remembers_count | last_remembered? | next_repetition = (last_remembered + repetition_interval(successful_remembers_count))
}

func TestExecuteQuizCommandWhenNoWordsToRepeat(t *testing.T) {

}
