package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractAddCommandWithParams(t *testing.T) {
	var expectedUserId int32 = 7
	expectedPhrase := "hello"
	expectedTranslation := "здравствуйте, привет"
	command, ok := extractCommandFrom(expectedUserId, "/add hello - здравствуйте, привет").(AddDictionaryEntryCommand)

	require.True(t, ok, "Command not extracted")
	require.Equal(t, expectedUserId, command.userId)
	require.Equal(t, expectedPhrase, command.phrase)
	require.Equal(t, expectedTranslation, command.translation)
}

func TestExtractQuizCommand(t *testing.T) {
	var expectedUserId int32 = 7
	command, ok := extractCommandFrom(expectedUserId, "/quiz").(QuizCommand)

	require.True(t, ok, "Command not extracted")
	require.Equal(t, expectedUserId, command.userId)
}
