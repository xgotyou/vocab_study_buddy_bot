package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// anonymous import for registering driver
	"github.com/ydb-platform/ydb-go-sdk/v3"
	yc "github.com/ydb-platform/ydb-go-yc"
)

var db *sql.DB

type DictionaryEntry struct {
	id                          int32
	userId                      int32
	phrase                      string
	translation                 string
	phraseMemorizationStep      uint16
	translationMemorizationStep uint16
	nextPhraseRepetition        time.Time
	nextTranslationRepetition   time.Time
}

func FindDictionaryEntryByUserIdAndPhrase(userId int32, phrase string) (DictionaryEntry, error) {
	row := db.QueryRow(fmt.Sprintf(
		"SELECT id, user_id, phrase, translation FROM dictionary_entries WHERE user_id = %v AND phrase = %q",
		userId,
		phrase))

	var entry DictionaryEntry
	err := row.Scan(&entry.id, &entry.userId, &entry.phrase, &entry.translation)

	if err != nil {
		return DictionaryEntry{}, err
	}

	return entry, err
}

func (entry *DictionaryEntry) Save() error {
	err := entry.validate()
	if err != nil {
		return err
	}

	q := fmt.Sprintf(
		`INSERT INTO dictionary_entries (user_id, phrase, translation, phrase_memorization_step, translation_memorization_step, 
			next_phrase_repetition, next_translation_repetition) 
		VALUES (%d, %q, %q, %d, %d, DateTime(%q), DateTime(%q));
		SELECT MAX(id) FROM dictionary_entries;`,
		entry.userId,
		entry.phrase,
		entry.translation,
		entry.phraseMemorizationStep,
		entry.translationMemorizationStep,
		entry.nextPhraseRepetition.Format(time.RFC3339),
		entry.nextTranslationRepetition.Format(time.RFC3339),
	)

	row := db.QueryRow(q)

	err = row.Scan(&entry.id)

	return err
}

func FindDictionaryEntryToRepeat(userId int32) (DictionaryEntry, error) {
	row := db.QueryRow(fmt.Sprintf(`
		SELECT id, user_id, phrase, translation, phrase_memorization_step, translation_memorization_step, 
			next_phrase_repetition, next_translation_repetition
		FROM dictionary_entries
		WHERE user_id = %d
		ORDER BY next_phrase_repetition ASC
		LIMIT 1`, userId))

	var entry DictionaryEntry
	err := row.Scan(&entry.id, &entry.userId, &entry.phrase, &entry.translation, &entry.phraseMemorizationStep,
		&entry.translationMemorizationStep, &entry.nextPhraseRepetition, &entry.nextTranslationRepetition)

	return entry, err
}

func (entry *DictionaryEntry) validate() error {
	if entry.userId == 0 {
		return fmt.Errorf("User ID can't be null")
	}

	if entry.phrase == "" {
		return fmt.Errorf("Phrase can't be blank")
	}

	if entry.translation == "" {
		return fmt.Errorf("Translation can't be blank")
	}

	return nil
}

func InitDB(ctx context.Context) {
	dsn := os.Getenv("YDB_CONNECTION_STRING")
	var ydbCredentials ydb.Option

	if strings.Contains(dsn, "localhost") {
		// Consider working in local developmant environment
		ydbCredentials = ydb.WithAnonymousCredentials()
	} else {
		// Consider working as Yandex Cloud Function
		ydbCredentials = yc.WithMetadataCredentials()
		log.Println("Using YC Metadata credentials")
	}

	nativeDriver := ydb.MustOpen(ctx, dsn, ydbCredentials)
	db = sql.OpenDB(ydb.MustConnector(nativeDriver))
}
