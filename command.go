package main

import "regexp"

type Command interface {
	Execute() (string, error)
}

func extractCommandFrom(userId int32, str string) Command {
	re := regexp.MustCompile(`^(/[^ ]*)(?: *([^\-]+[^ ]) *- *(.+))?`)

	matches := re.FindStringSubmatch(str)

	if matches == nil {
		return AddDictionaryEntryCommand{}
	}

	switch matches[1] {
	case "/add":
		return AddDictionaryEntryCommand{
			userId:      userId,
			phrase:      matches[2],
			translation: matches[3],
		}
	case "/quiz":
		return QuizCommand{userId: userId}
	}

	return AddDictionaryEntryCommand{}
}
