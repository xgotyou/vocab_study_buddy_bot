package main

import (
	"testing"
)

func TestExecuteAddDictionaryEntryCommand(t *testing.T) {
	command := AddDictionaryEntryCommand{
		userId:      7,
		phrase:      "hi",
		translation: "привет",
	}

	expectedResult := `Successfully added entry for "hi"`

	actualResult, err := command.Execute()

	if err != nil {
		t.Fatalf("Expected error to be nil, got %q", err.Error())
	}

	if actualResult != expectedResult {
		t.Fatalf("Expected result to be %q, got %q", expectedResult, actualResult)
	}

	entry, err := FindDictionaryEntryByUserIdAndPhrase(command.userId, command.phrase)

	if err != nil {
		t.Fatalf("Error during attempt to find added dictionary entry: %q", err.Error())
	}

	if entry.translation != command.translation {
		t.Fatalf("Entry translation is expected to be %q, was %q", command.translation, entry.translation)
	}
}
