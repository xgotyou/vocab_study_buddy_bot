package main

import (
	"reflect"
	"testing"
)

func TestSaveNew(t *testing.T) {
	entry := DictionaryEntry{
		userId:      7,
		phrase:      "phrase",
		translation: "translation",
	}

	err := entry.Save()

	if err != nil {
		t.Fatalf("Error expected to be nil, but was %q", err.Error())
	}

	if entry.id == 0 {
		t.Fatalf("Entry is expected to get non null ID after persistance")
	}
}

func TestSaveInvalidUserId(t *testing.T) {
	entry := DictionaryEntry{
		phrase:      "phrase",
		translation: "translation",
	}

	expectedError := "User ID can't be null"
	err := entry.Save()

	if err == nil {
		t.Fatalf("Error is expected to not be nil")
	}

	if err.Error() != expectedError {
		t.Fatalf("Expect error to be %q, got %q", expectedError, err.Error())
	}
}

func TestSaveBlankPhrase(t *testing.T) {
	entry := DictionaryEntry{
		userId:      7,
		translation: "translation",
	}

	expectedError := "Phrase can't be blank"
	err := entry.Save()

	if err == nil {
		t.Fatalf("Error is expected to not be nil")
	}

	if err.Error() != expectedError {
		t.Fatalf("Expect error to be %q, got %q", expectedError, err.Error())
	}
}

func TestSaveBlankTranslation(t *testing.T) {
	entry := DictionaryEntry{
		userId: 7,
		phrase: "phrase",
	}

	expectedError := "Translation can't be blank"
	err := entry.Save()

	if err == nil {
		t.Fatalf("Error is expected to not be nil")
	}

	if err.Error() != expectedError {
		t.Fatalf("Expect error to be %q, got %q", expectedError, err.Error())
	}
}

func TestFindDictionaryEntryByUserIdAndPhraseWithValidEntry(t *testing.T) {
	clearDB()

	expectedEntry := DictionaryEntry{
		userId:      7,
		phrase:      "phrase",
		translation: "translation",
	}

	err := expectedEntry.Save()
	if err != nil {
		t.Fatalf("Error while trying to create entry")
	}

	actualEntry, err := FindDictionaryEntryByUserIdAndPhrase(expectedEntry.userId, expectedEntry.phrase)

	if err != nil {
		t.Fatalf("Expect error to be nil, got %q", err)
	}

	if !reflect.DeepEqual(expectedEntry, actualEntry) {
		t.Fatalf("Expect dictionary entry to be %v, got %v", expectedEntry, actualEntry)
	}
}

func TestFindDictionaryEntryByUserIdAndPhraseWithNonExistingPhrase(t *testing.T) {
	existingEntry := DictionaryEntry{
		userId:      7,
		phrase:      "phrase",
		translation: "translation",
	}

	err := existingEntry.Save()
	if err != nil {
		t.Fatalf("Error while trying to create entry")
	}

	_, err = FindDictionaryEntryByUserIdAndPhrase(existingEntry.userId, "unknown phrase")

	if err == nil {
		t.Fatalf("Error is expected to not be nil")
	}
}

func TestFindDictionaryEntryByUserIdAndPhraseWithNonExistingUser(t *testing.T) {
	existingEntry := DictionaryEntry{
		userId:      7,
		phrase:      "phrase",
		translation: "translation",
	}

	err := existingEntry.Save()
	if err != nil {
		t.Fatalf("Error while trying to create entry")
	}

	_, err = FindDictionaryEntryByUserIdAndPhrase(1, existingEntry.phrase)

	if err == nil {
		t.Fatalf("Error is expected to not be nil")
	}
}
