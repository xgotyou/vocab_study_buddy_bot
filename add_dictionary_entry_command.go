package main

import "fmt"

type AddDictionaryEntryCommand struct {
	userId      int32
	phrase      string
	translation string
}

func (command AddDictionaryEntryCommand) Execute() (string, error) {
	entry := DictionaryEntry{
		userId:      command.userId,
		phrase:      command.phrase,
		translation: command.translation,
	}

	err := entry.Save()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully added entry for %q", command.phrase), nil
}
